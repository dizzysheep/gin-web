package ginc

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func Query(c *gin.Context, key string) string {
	if val := c.Param(key); val != "" {
		return val
	}

	return c.Query(key)
}

// GetPage returns paging parameter.
func GetPage(c *gin.Context) int {
	ret, _ := strconv.Atoi(c.Query("page"))
	if 1 > ret {
		ret = 1
	}

	return ret
}

// GetPageSize returns paging parameter.
func GetPageSize(c *gin.Context) int {
	ret, _ := strconv.Atoi(c.Query("size"))
	if 1 > ret {
		ret = 10
	}

	return ret
}

// GetString returns the input value by key string or the default value while it's present and input is blank
func GetString(c *gin.Context, key string, def ...string) string {
	if v := Query(c, key); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetInt returns input as an int or the default value while it's present and input is blank
func GetInt(c *gin.Context, key string, def ...int) (i int) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		i = def[0]
		return
	}
	i, _ = strconv.Atoi(strv)
	return
}

// GetInt8 return input as an int8 or the default value while it's present and input is blank
func GetInt8(c *gin.Context, key string, def ...int8) (i8 int8) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		i8 = def[0]
		return
	}
	i64, _ := strconv.ParseInt(strv, 10, 8)
	i8 = int8(i64)
	return
}

// GetUint8 return input as an uint8 or the default value while it's present and input is blank
func GetUint8(c *gin.Context, key string, def ...uint8) (u8 uint8) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		u8 = def[0]
		return
	}
	u64, _ := strconv.ParseUint(strv, 10, 8)
	u8 = uint8(u64)
	return
}

// GetInt16 returns input as an int16 or the default value while it's present and input is blank
func GetInt16(c *gin.Context, key string, def ...int16) (i16 int16) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		i16 = def[0]
		return
	}
	i64, _ := strconv.ParseInt(strv, 10, 16)
	i16 = int16(i64)
	return
}

// GetUint16 returns input as an uint16 or the default value while it's present and input is blank
func GetUint16(c *gin.Context, key string, def ...uint16) (u16 uint16) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		u16 = def[0]
		return
	}
	u64, _ := strconv.ParseUint(strv, 10, 16)
	u16 = uint16(u64)
	return
}

// GetInt32 returns input as an int32 or the default value while it's present and input is blank
func GetInt32(c *gin.Context, key string, def ...int32) (i32 int32) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		i32 = def[0]
		return
	}
	i64, _ := strconv.ParseInt(strv, 10, 32)
	i32 = int32(i64)
	return
}

// GetUint32 returns input as an uint32 or the default value while it's present and input is blank
func GetUint32(c *gin.Context, key string, def ...uint32) (u32 uint32) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		u32 = def[0]
		return
	}
	u64, _ := strconv.ParseUint(strv, 10, 32)
	u32 = uint32(u64)
	return
}

// GetInt64 returns input value as int64 or the default value while it's present and input is blank.
func GetInt64(c *gin.Context, key string, def ...int64) (i64 int64) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		i64 = def[0]
		return
	}
	i64, _ = strconv.ParseInt(strv, 10, 64)
	return
}

// GetUint64 returns input value as uint64 or the default value while it's present and input is blank.
func GetUint64(c *gin.Context, key string, def ...uint64) (u64 uint64) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		u64 = def[0]
		return
	}
	u64, _ = strconv.ParseUint(strv, 10, 64)
	return
}

// GetBool returns input value as bool or the default value while it's present and input is blank.
func GetBool(c *gin.Context, key string, def ...bool) (b bool) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		b = def[0]
		return
	}
	b, _ = strconv.ParseBool(strv)
	return
}

// GetFloat returns input value as float64 or the default value while it's present and input is blank.
func GetFloat(c *gin.Context, key string, def ...float64) (f64 float64) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		f64 = def[0]
		return
	}

	f64, _ = strconv.ParseFloat(strv, 64)
	return
}

func GetIntE(c *gin.Context, key string, def ...int) (int, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(strv)
}

func GetInt8E(c *gin.Context, key string, def ...int8) (int8, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 8)
	return int8(i64), err
}

func GetUint8E(c *gin.Context, key string, def ...uint8) (uint8, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 8)
	return uint8(u64), err
}

func GetInt16E(c *gin.Context, key string, def ...int16) (int16, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 16)
	return int16(i64), err
}

func GetUint16E(c *gin.Context, key string, def ...uint16) (uint16, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 16)
	return uint16(u64), err
}

func GetInt32E(c *gin.Context, key string, def ...int32) (int32, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 32)
	return int32(i64), err
}

func GetUint32E(c *gin.Context, key string, def ...uint32) (uint32, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 32)
	return uint32(u64), err
}

func GetInt64E(c *gin.Context, key string, def ...int64) (int64, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseInt(strv, 10, 64)
}

func GetUint64E(c *gin.Context, key string, def ...uint64) (uint64, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseUint(strv, 10, 64)
}

func GetBoolE(c *gin.Context, key string, def ...bool) (bool, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseBool(strv)
}

func GetFloatE(c *gin.Context, key string, def ...float64) (float64, error) {
	strv := Query(c, key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseFloat(strv, 64)
}
