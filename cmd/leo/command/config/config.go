package config

import (
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/fatih/color"
	"github.com/fatih/structtag"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type Config struct {
	DBType                              string `json:"db_type" default:"mysql" description:"数据库类型，可以是mysql或sqlite3"`
	DBAddr                              string `json:"db_addr" default:"yitu:password@127.0.0.1/video_mario" description:"MySQL数据库地址"`
	DBMaxIdles                          int    `json:"db_max_idles" default:"5"`
	DBMaxConns                          int    `json:"db_max_conns" default:"10"`
	HttpAddr                            string `json:"http_addr" default:":8989" description:"Video_mario绑定的HTTP服务地址"`
	HttpReadTimeout                     int    `json:"http_read_timeout" default:"10"`
	HttpWriteTimeout                    int    `json:"http_write_timeout" default:"20"`
	LogDir                              string `json:"log_dir" default:"./log/" description:"日志目录"`
	LogLevel                            string `json:"log_level" default:"info" description:"日志等级"`
	LogRotationSize                     int    `json:"log_rotation_size" default:"100"`
	LogExpireTime                       int    `json:"log_expire_time" default:"7"`
	LogTotalSize                        int    `json:"log_total_size" default:"1024"`
	AiType                              string `json:"ai_type" default:"fp" description:"启用哪种计算服务，可以是cube，fp或memory"`
	CubeAddr                            string `json:"cube_addr" default:"127.0.0.1:8091" description:"魔方计算服务地址"`
	CubeSecurityEnabled                 bool   `json:"cube_security_enabled" default:"false"`
	CubeUser                            string `json:"cube_user" default:"admin"`
	CubePass                            string `json:"cube_pass" default:"admin123"`
	FPAddr                              string `json:"fp_addr" default:"test123:test123456@127.0.0.1" description:"FP计算服务地址"`
	FeatureVersion                      int    `json:"feature_version" default:"2230" description:"人脸模型Feature Version"`
	SimilarQueueSize                    int    `json:"similar_queue_size" default:"100"`
	SimilarTaskSize                     int    `json:"similar_task_size" default:"100"`
	RegisterChanSize                    int    `json:"register_chan_size" default:"1000"`
	GuestChanSize                       int    `json:"guest_chan_size" default:"1000"`
	ExtraChanSize                       int    `json:"extra_chan_size" default:"1000"`
	ExtraPersonChanSize                 int    `json:"extra_person_chan_size" default:"1000"`
	StatsChanSize                       int    `json:"stats_chan_size" default:"1000"`
	Parallel1v1Count                    int    `json:"parallel_1v1_count" default:"10"`
	Parallel1vnCount                    int    `json:"parallel_1vn_count" default:"5"`
	Time1v1Threshold                    int    `json:"time_1v1_threshold" default:"10"`
	BatchDBCount                        int    `json:"batch_db_count" default:"1000"`
	SaveFeature                         bool   `json:"save_feature" default:"true" description:"是否保存feature"`
	RetrievalStrategy                   string `json:"retrieval_strategy" default:"simple"`
	AiClientRetryTimes                  int    `json:"ai_client_retry_times" default:"3"`
	AiClientRetryEnabled                bool   `json:"ai_client_retry_enabled" default:"true"`
	CompanionSecondThreshold            int    `json:"companion_second_threshold" default:"3"`
	CompanionDeduplicateSecondThreshold int    `json:"companion_deduplicate_second_threshold" default:"0"`
	ReenterSecondThreshold              int    `json:"reenter_second_threshold" default:"10"`
	StructureCameraName                 string `json:"structure_camera_name" default:"_GetMetaDetail"`
	UpdateStructureInfo                 bool   `json:"update_structure_info" default:"false" description:"是否更新结构化信息"`
	ExcludeHat                          bool   `json:"exclude_hat" default:"false"`
	ExcludeMask                         bool   `json:"exclude_mask" default:"false"`
	ExcludeSunglass                     bool   `json:"exclude_sunglass" default:"false"`
	EnterAreaCheck                      bool   `json:"enter_area_check" default:"true" description:"是否开启进入区域检测"`
	ExitAreaCheck                       bool   `json:"exit_area_check" default:"true" description:"是否开启离开区域检测"`
	NotifyUrl                           string `json:"notify_url" default:"http://127.0.0.1:8095/vbox/v1/event"`
	NotifyEnabled                       bool   `json:"notify_enabled" default:"false" description:"是否开启通知推送"`
	ImageServer                         string `json:"image_server" default:"127.0.0.1:8092"`
	Debug                               bool   `json:"debug" default:"false" description:"Debug模式"`
	MaxPersonFeatures                   int    `json:"max_person_features" default:"0" description:"路人库中每人最大记录数"`
	DataReservedDays                    int    `json:"data_reserved_days" default:"30" description:"清理数据的时候保留多少天数据, 0表示全清"`
}

var Current Config

func InitFlag(rootCmd *cobra.Command) {
	t := reflect.TypeOf(Current)
	v := reflect.ValueOf(&Current)
	for i := 0; i < t.NumField(); i++ {
		tags, _ := structtag.Parse(string(t.Field(i).Tag))
		jsonTag, _ := tags.Get("json")
		if jsonTag == nil {
			panic("config field must have json tag")
		}
		defaultTag, _ := tags.Get("default")
		if defaultTag == nil {
			panic("config field must have default tag")
		}
		descriptionTag, _ := tags.Get("description")
		switch t.Field(i).Type {
		case reflect.TypeOf(""):
			if descriptionTag != nil {
				rootCmd.PersistentFlags().StringVar(v.Elem().Field(i).Addr().Interface().(*string), jsonTag.Name, defaultTag.Name, descriptionTag.Name)
			} else {
				v.Elem().Field(i).SetString(defaultTag.Name)
			}
		case reflect.TypeOf(0):
			name, _ := strconv.Atoi(defaultTag.Name)
			if descriptionTag != nil {
				rootCmd.PersistentFlags().IntVar(v.Elem().Field(i).Addr().Interface().(*int), jsonTag.Name, name, descriptionTag.Name)
			} else {
				v.Elem().Field(i).SetInt(int64(name))
			}
		case reflect.TypeOf(true):
			name, _ := strconv.ParseBool(defaultTag.Name)
			if descriptionTag != nil {
				rootCmd.PersistentFlags().BoolVar(v.Elem().Field(i).Addr().Interface().(*bool), jsonTag.Name, name, descriptionTag.Name)
			} else {
				v.Elem().Field(i).SetBool(name)
			}
		}
	}
}

func Init(configFile string) {
	if len(configFile) == 0 {
		return
	}
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		color.Red("Cannot read specified configuration file %s: %v", configFile, err)
		return
	}
	fileJson := string(bytes)
	t := reflect.TypeOf(Current)
	v := reflect.ValueOf(&Current)
	for i := 0; i < t.NumField(); i++ {
		tags, _ := structtag.Parse(string(t.Field(i).Tag))
		jsonTag, _ := tags.Get("json")
		if jsonTag == nil {
			panic("config field must have json tag")
		}
		if gjson.Get(fileJson, jsonTag.Name).Exists() {
			switch t.Field(i).Type {
			case reflect.TypeOf(""):
				v.Elem().Field(i).SetString(gjson.Get(fileJson, jsonTag.Name).String())
			case reflect.TypeOf(0):
				v.Elem().Field(i).SetInt(gjson.Get(fileJson, jsonTag.Name).Int())
			case reflect.TypeOf(true):
				v.Elem().Field(i).SetBool(gjson.Get(fileJson, jsonTag.Name).Bool())
			}
		}
	}
}

func Update(key string, value string) error {
	t := reflect.TypeOf(Current)
	v := reflect.ValueOf(&Current)
	found := false
	for i := 0; i < t.NumField(); i++ {
		tags, _ := structtag.Parse(string(t.Field(i).Tag))
		jsonTag, _ := tags.Get("json")
		if jsonTag.Name == key {
			found = true
			switch t.Field(i).Type {
			case reflect.TypeOf(""):
				v.Elem().Field(i).SetString(value)
			case reflect.TypeOf(0):
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				v.Elem().Field(i).SetInt(intValue)
			case reflect.TypeOf(true):
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					return err
				}
				v.Elem().Field(i).SetBool(boolValue)
			}
		}
	}
	if !found {
		return errors.New("config not exist: " + key)
	}
	return nil
}
