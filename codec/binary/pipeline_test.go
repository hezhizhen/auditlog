package binary_test

import (
	"io"
	"sync"
	"testing"

	"github.com/containerssh/auditlog/codec"
	"github.com/containerssh/auditlog/codec/binary"
	"github.com/containerssh/auditlog/message"

	"github.com/stretchr/testify/assert"
)

func createPipeline() (codec.Encoder, codec.Decoder) {
	return binary.NewEncoder(), binary.NewDecoder()
}

func testPipeline(t *testing.T, msg *message.Message) {
	encoder, decoder := createPipeline()

	pipeReader, pipeWriter := io.Pipe()

	messageChannel := make(chan message.Message, 1)
	messageChannel <- *msg
	close(messageChannel)

	storage := codec.NewStorageWriterProxy(pipeWriter)

	wg := sync.WaitGroup{}
	wg.Add(2)
	var encodeError error = nil
	var decodeError error = nil
	var decodedMessage *message.Message
	go func() {
		defer wg.Done()
		err := encoder.Encode(messageChannel, storage)
		if err != nil {
			encodeError = err
		}
	}()
	go func() {
		defer wg.Done()
		decodedMessageChannel, errorsChannel, done := decoder.Decode(pipeReader)
		select {
		case decodedMessage = <-decodedMessageChannel:
		case decodeError = <-errorsChannel:
		}
		<-done
	}()
	wg.Wait()
	if encodeError != nil {
		assert.Failf(t, "encoding error", "encoding error (%w)", encodeError)
	}
	if decodeError != nil {
		assert.Failf(t, "decoding error", "decoding error (%w)", decodeError)
	}

	assert.True(t, msg.Equals(decodedMessage))
}

func TestTypeConnect(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeConnect,
		Payload: message.PayloadConnect{
			RemoteAddr: "127.0.0.1",
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeDisconnect(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeDisconnect,
		Payload:      nil,
		ChannelID:    -1,
	}

	testPipeline(t, &msg)
}

func TestTypeAuthPassword(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeAuthPassword,
		Payload: message.PayloadAuthPassword{
			Username: "foo",
			Password: []byte("bar"),
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeAuthPasswordSuccessful(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeAuthPasswordSuccessful,
		Payload: message.PayloadAuthPassword{
			Username: "foo",
			Password: []byte("bar"),
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeAuthPasswordFailed(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeAuthPasswordFailed,
		Payload: message.PayloadAuthPassword{
			Username: "foo",
			Password: []byte("bar"),
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeAuthPasswordBackendError(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeAuthPasswordBackendError,
		Payload: message.PayloadAuthPasswordBackendError{
			Username: "foo",
			Password: []byte("bar"),
			Reason:   "test",
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeAuthPubKey(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeAuthPubKey,
		Payload: message.PayloadAuthPubKey{
			Username: "foo",
			Key:      []byte("bar"),
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeAuthPubKeySuccessful(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeAuthPubKeySuccessful,
		Payload: message.PayloadAuthPubKey{
			Username: "foo",
			Key:      []byte("bar"),
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeAuthPubKeyFailed(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeAuthPubKeyFailed,
		Payload: message.PayloadAuthPubKey{
			Username: "foo",
			Key:      []byte("bar"),
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeAuthPubKeyBackendError(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeAuthPubKeyBackendError,
		Payload: message.PayloadAuthPubKeyBackendError{
			Username: "foo",
			Key:      []byte("bar"),
			Reason:   "test",
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeGlobalRequestUnknown(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeGlobalRequestUnknown,
		Payload: message.PayloadGlobalRequestUnknown{
			RequestType: "channel",
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeNewChannel(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeNewChannel,
		Payload: message.PayloadNewChannel{
			ChannelType: "session",
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeNewChannelFailed(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeNewChannelFailed,
		Payload: message.PayloadNewChannelFailed{
			ChannelType: "session",
			Reason:      "test",
		},
		ChannelID: -1,
	}

	testPipeline(t, &msg)
}

func TestTypeNewChannelSuccessful(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeNewChannelSuccessful,
		Payload: message.PayloadNewChannelSuccessful{
			ChannelType: "session",
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestUnknownType(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestUnknownType,
		Payload: message.PayloadChannelRequestUnknownType{
			RequestType: "test",
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestDecodeFailed(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestDecodeFailed,
		Payload: message.PayloadChannelRequestDecodeFailed{
			RequestType: "test",
			Reason:      "test",
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestSetEnv(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestSetEnv,
		Payload: message.PayloadChannelRequestSetEnv{
			Name:  "foo",
			Value: "bar",
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestExec(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestExec,
		Payload: message.PayloadChannelRequestExec{
			Program: "bash",
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestPty(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestPty,
		Payload: message.PayloadChannelRequestPty{
			Columns: 80,
			Rows:    25,
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestShell(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestShell,
		Payload:      message.PayloadChannelRequestShell{},
		ChannelID:    0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestSignal(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestSignal,
		Payload: message.PayloadChannelRequestSignal{
			Signal: "TERM",
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestSubsystem(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestSubsystem,
		Payload: message.PayloadChannelRequestSubsystem{
			Subsystem: "sftp",
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}

func TestTypeChannelRequestWindow(t *testing.T) {
	msg := message.Message{
		ConnectionID: []byte("1234"),
		Timestamp:    1234,
		MessageType:  message.TypeChannelRequestWindow,
		Payload: message.PayloadChannelRequestWindow{
			Columns: 80,
			Rows:    25,
		},
		ChannelID: 0,
	}

	testPipeline(t, &msg)
}
