package util
import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chilts/sid"
	"github.com/kjk/betterguid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"
)

func genXid() {
	id := xid.New()
	fmt.Printf("github.com/rs/xid:           %s\n", id.String())
}

func genKsuid() {
	id := ksuid.New()
	fmt.Printf("github.com/segmentio/ksuid:  %s\n", id.String())
}

func genBetterGUID() {
	id := betterguid.New()
	fmt.Printf("github.com/kjk/betterguid:   %s\n", id)
}

// 20个字母
func genUlid() {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	fmt.Printf("github.com/oklog/ulid:       %s\n", id.String())
}

func genSonyflake() {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	// Note: this is base16, could shorten by encoding as base62 string
	fmt.Printf("github.com/sony/sonyflake:   %x\n", id)
}

func genSid() {
	id := sid.Id()
	fmt.Printf("github.com/chilts/sid:       %s\n", id)
}

func genUUIDv4() {
	id := uuid.NewV4()
	// if err != nil {
	// 	log.Fatalf("uuid.NewV4() failed with %s\n", err)
	// }
	fmt.Printf("github.com/satori/go.uuid:   %s\n", id)
}

// 会话ID
func GenerateTalkSessionID() (string){
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

// 会话ID
func GenerateTalkSubsessionID() (string){
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

// APP ID
func GenerateAppID() string {
	id := xid.New()
	return id.String()
}

// 用户ID
func GenerateUserID() string {
	id := xid.New()
	return id.String()
}

// 房间ID
func GenerateRoomID() string {
	id := xid.New()
	return id.String()
}

// 用户操作记录ID
func GenerateUserOperateId() string {
	id := xid.New()
	return id.String()
}

// 用户信息记录ID
func GenerateUserInfoId() string {
	id := xid.New()
	return id.String()
}



//func Generate20BitUuid()  (string){
//	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
//	id, err := flake.NextID()
//	if err != nil {
//		log.Fatalf("flake.NextID() failed with %s\n", err)
//	}
//	return id
//}