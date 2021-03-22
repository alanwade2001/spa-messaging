package services

import (
	"log"
	"reflect"

	"github.com/alanwade2001/spa-messaging/spa-msg-initiation-instruction/generated/initiation"
	"github.com/alanwade2001/spa-messaging/spa-msg-initiation-instruction/types"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"k8s.io/klog/v2"

	mgo "github.com/alanwade2001/spa-common/mongo"
)

type Message struct {
}

func NewMessage() types.MessageAPI {
	return &Message{}
}

func (m *Message) Process(body []byte) error {
	model := new(initiation.InitiationModel)
	if err := model.UnmarshalJSON(body); err != nil {
		return err
	}

	log.Printf("%+v", model)

	structcodec, _ := bsoncodec.NewStructCodec(bsoncodec.JSONFallbackStructTagParser)
	reg := bson.NewRegistryBuilder().
		RegisterTypeEncoder(reflect.TypeOf(initiation.InitiationModel{}), structcodec).
		RegisterTypeDecoder(reflect.TypeOf(initiation.InitiationModel{}), structcodec).Build()

	uriTemplate := viper.GetViper().GetString("MONGODB_URI_TEMPLATE")
	mongoUser := viper.GetViper().GetString("MONGODB_USER")
	mongoPassword := viper.GetViper().GetString("MONGODB_PASSWORD")
	mongoDatabase := viper.GetViper().GetString("MONGODB_DATABASE")
	mongoCollection := viper.GetViper().GetString("MONGODB_COLLECTION")
	mongoTimeout := viper.GetViper().GetDuration("MONGODB_TIMEOUT")

	mongoService := mgo.NewMongoService(uriTemplate, mongoUser, mongoPassword, mongoDatabase, mongoCollection, mongoTimeout, reg)
	conn := mongoService.Connect()
	defer conn.Disconnect()

	result, err := mongoService.GetCollection(conn).InsertOne(conn.Ctx, model)

	if err != nil {
		klog.Warningf("Could not create Initiation: %v", err)
		return err
	}

	klog.Infof("result:[%+v]", result)
	klog.Infof("initiation:[%+v]", model)

	return nil
}
