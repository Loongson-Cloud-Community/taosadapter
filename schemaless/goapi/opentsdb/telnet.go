package opentsdb

import (
	"crypto/md5"
	"fmt"
	"unsafe"

	"github.com/taosdata/taosadapter/schemaless"
	"github.com/taosdata/taosadapter/schemaless/goapi/parser/opentsdb/telnet"
	"github.com/taosdata/taosadapter/tools/pool"
)

func InsertOpentsdbTelnet(taosConnect unsafe.Pointer, data, db string) error {
	executor, err := schemaless.NewExecutor(taosConnect)
	if err != nil {
		return err
	}
	point, err := telnet.Unmarshal(data)
	if err != nil {
		return err
	}
	defer telnet.CleanPoint(point)
	tagNames := make([]string, len(point.Tags))
	tagValues := make([]string, len(point.Tags))

	b := pool.BytesPoolGet()

	b.WriteByte('`')
	b.WriteString(point.Metric)
	b.WriteByte('`')
	stableName := b.String()
	b.Reset()
	b.WriteString(stableName)
	tmpB := pool.BytesPoolGet()
	for i, tag := range point.Tags {
		tmpB.Reset()

		tmpB.WriteByte('`')
		tmpB.WriteString(tag.Key)
		tmpB.WriteByte('`')
		name := tmpB.String()
		tagNames[i] = name
		tagValues[i] = tag.Value
		b.WriteByte(',')
		b.WriteString(name)
		b.WriteByte('=')
		b.WriteString(tag.Value)
	}
	pool.BytesPoolPut(tmpB)
	tableName := fmt.Sprintf("t_%x", md5.Sum(b.Bytes()))
	sql, err := executor.InsertTDengine(&schemaless.InsertLine{
		DB:         db,
		Ts:         point.Ts,
		TableName:  tableName,
		STableName: stableName,
		Fields: map[string]interface{}{
			valueField: point.Value,
		},
		TagNames:  tagNames,
		TagValues: tagValues,
	})
	if err != nil {
		logger.WithError(err).Errorln(sql)
		return err
	}
	return nil
}
