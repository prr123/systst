package posx_mq

import (
	"bytes"
	"log"
	"os"
	"os/exec"
//	"strconv"
	"strings"
	"syscall"
	"errors"
)

//var cadxQueue asset.CadxQueue

// Represents the message queue
type MessageQueue struct {
	handler int
	name    string
	recvBuf *receiveBuffer
//	recvBuf []byte
}

type QueueStats struct {
	maxMessages         int64
	curMessages         int64
	queueUsedPercentage float64
}

// Represents the message queue attribute
type MessageQueueAttribute struct {
	flags   int
	maxMsg  int
	msgSize int
	curMsgs int
}

type Posix_msg_que struct {
	Name string
	MaxMessages int
	MaxMessageSize int
}

// NewMessageQueue returns an instance of the message queue given a QueueConfig.
// Create the message queue with some default settings.
// ma is a struct:   struct mq_attr ma;      // message queue attributes
//    mq = mq_open("/test_queue", O_RDWR | O_CREAT, 0700, &ma);

func NewMessageQueue(name string, flag, permission, max_msgs, maxmsgsize int) (*MessageQueue, error) {
// change: don't need
//	cadxQueue = asset.GetCadxQueueFilterDefinition(filterConfigRootPath)

	var mq MessageQueue

	res := strings.Contains(name,"/")
	if res {
		err := errors.New("error: msg queue name contains /")
		return nil, err
	}

	name = "/" + name
//	fmt.Println("queue name: ",name)

// change
	unlinkQueue(name)

	h, err := mq_open(name, flag, permission, max_msgs, maxmsgsize)

	if err != nil {
		return nil, err
	}

// change
	recvBuf, err := newReceiveBuffer(int(maxmsgsize))
	if err != nil {
		mq_close(h)
		return nil, err
	}

	mq.handler = h
	mq.name = name
	mq.recvBuf = recvBuf

	return &mq, nil
}

// change
// get new without unlink

func GetMessageQueue(name string, flag, permission, max_msgs, maxmsgsize int) (*MessageQueue, error) {

	var mq MessageQueue

	res := strings.Contains(name,"/")
	if res {
		err := errors.New("error: msg queue name contains /")
		return nil, err
	}

	name = "/" + name

//	fmt.Println("queue name: ",name)
// change
//	unlinkQueue(name)

	
	h, err := mq_open(name, flag, permission, max_msgs, maxmsgsize)

	if err != nil {
		return nil, err
	}

// change
	recvBuf, err := newReceiveBuffer(int(maxmsgsize))
	if err != nil {
		mq_close(h)
		return nil, err
	}

	mq.handler = h
	mq.name = name
	mq.recvBuf = recvBuf

	return &mq, nil
}


// not necessary kernel has been stable since 3.5
func getKernelVersion() string {
	cmd := exec.Command("uname", "-r")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		log.Fatalf("I! error while executing uname command to get kernel version %s", err.Error())
	}
	kernelRelease := strings.TrimSpace(string(cmdOutput.Bytes()))
	return strings.Split(kernelRelease, "-")[0]
}

// Send sends message to the message queue.
func (mq *MessageQueue) Send(data []byte, priority uint) error {
	_, err := mq_send(mq.handler, data, priority)
	return err
}

// Receive receives message from the message queue.
// *** review recvbuf
func (mq *MessageQueue) Receive() ([]byte, uint, error) {
	return mq_receive(mq.handler, mq.recvBuf)
}

// TimedReceive receives message from the POSIX queue on a timely basis
// need to modify time
func (mq *MessageQueue) TimedReceive() ([]byte, uint, error) {
	return mq_timedreceive(mq.handler, mq.recvBuf)
}

// FIXME Don't work because of signal portability.
// Notify set signal notification to handle new messages.
// *** review
func (mq *MessageQueue) Notify(sigNo syscall.Signal) error {
	_, err := mq_notify(mq.handler, int(sigNo))
	return err
}

// Close closes the message queue.
func (mq *MessageQueue) Close() error {
	mq.recvBuf.free()

	_, err := mq_close(mq.handler)
	return err
}

// Unlink deletes the message queue.
func (mq *MessageQueue) Unlink() error {
	mq.Close()

	_, err := mq_unlink(mq.name)
	return err
}

// changed prr
func unlinkQueue(name string) error {
	_, err := mq_unlink(name)
	return err
}

// *** need to change
func changeQueuePermission(queueName string, posixQueueMountPath string) {
	var queuePath strings.Builder
	queuePath.WriteString(posixQueueMountPath)
	queuePath.WriteString(queueName)
	err := os.Chmod(queuePath.String(), 0442)
	if err != nil {
		log.Printf("in error while changing the permisssion for queue - %s , error : %s", queuePath, err.Error())
	}
}


