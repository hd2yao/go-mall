package util

import (
    "errors"
    "time"

    "github.com/jinzhu/copier"

    "github.com/hd2yao/go-mall/common/enum"
)

// CopyProperties 将 src(源对象) 的属性复制到 dst(目标对象)
func CopyProperties(src, dst interface{}) error {
    err := copier.CopyWithOption(dst, src, copier.Option{
        IgnoreEmpty: true, // 忽略空值，如果源对象的字段值为空，目标对象相应字段不会被覆盖
        DeepCopy:    true, // 深拷贝
        Converters: []copier.TypeConverter{ // 自定义类型转换器
            { // time.Time 转换成字符串
                SrcType: time.Time{},
                DstType: copier.String,
                Fn: func(src interface{}) (dst interface{}, err error) {
                    s, ok := src.(time.Time)
                    if !ok {
                        return nil, errors.New("src type is not time.Time")
                    }
                    return s.Format(enum.TimeFormatHyphenedYMDHIS), nil
                },
            },
            { // 字符串转换成 time.Time
                SrcType: copier.String,
                DstType: time.Time{},
                Fn: func(src interface{}) (dst interface{}, err error) {
                    s, ok := src.(string)
                    if !ok {
                        return nil, errors.New("src type is not time format string")
                    }

                    //// 匹配 YYYY-MM-DD HH:MM:SS 格式的字符串
                    //pattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$` // YYYY-MM-DD HH:MM:SS
                    //matched, _ := regexp.MatchString(pattern, s)
                    //if matched {
                    //	return time.Parse(enum.TimeFormatHyphenedYMDHIS, s)
                    //}
                    //return nil, errors.New("src type is not time format string")

                    // 直接用 time.Parse 解析，因为 time.Parse 本身会返回一个错误
                    t, err := time.Parse(enum.TimeFormatHyphenedYMDHIS, s)
                    if err != nil {
                        return nil, errors.New("failed to parse time string")
                    }
                    return t, nil
                },
            },
        },
    })

    return err
}
