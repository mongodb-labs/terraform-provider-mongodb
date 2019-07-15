package ssh

import (
	"bytes"
	"fmt"
	"github.com/MihaiBojin/terraform-provider-mongodb/mongodb/util"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/terraform/communicator/remote"
	tfssh "github.com/hashicorp/terraform/communicator/ssh"
	"github.com/hashicorp/terraform/terraform"
)

// Client abstracts all usage of SSH primitives into a nice, easy-to-use API
type Client struct {
	ephemeralState *terraform.EphemeralState
	communicator   *tfssh.Communicator
	mutex          *sync.Mutex
}

// NewClient creates a new SSH client
func NewClient(params ...func(*Connection) error) (*Client, error) {
	connInfo := &Connection{}
	for _, paramFunc := range params {
		err := paramFunc(connInfo)
		util.PanicOnNonNilErr(err)
	}

	state := connInfo.toEphemeralState()
	ephemeral := &terraform.InstanceState{ID: "ssh", Attributes: make(map[string]string), Meta: make(map[string]interface{}), Tainted: false, Ephemeral: *state}
	communicator, err := tfssh.New(ephemeral)
	if err != nil {
		return nil, fmt.Errorf("ssh.NewClient: could not initialize communicator, err=%v", err)
	}

	c := &Client{ephemeralState: state, communicator: communicator, mutex: &sync.Mutex{}}
	return c, c.connect()
}

func (c *Client) connect() error {
	// ensure we don't replace the communicator while another operation is in-progress, thus losing its output
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// connect to remote host
	err := c.communicator.Connect(WithLogging())

	if err != nil {
		res := Result{Cmd: "communicator.Connect", Err: err}
		return fmt.Errorf("ssh.connect: could not connect to remote host, err=%v", res)
	}

	return nil
}

// Upload uploads all contents from the specified io.Reader to the remote path
func (c *Client) UploadData(remotePath string, input io.Reader) (res Result) {
	// ensure we correctly retrieve the output associated with this command
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// print debugging info once the command completes
	defer res.Debug()

	log.Printf("Uploading data to the remote host at: %s", remotePath)
	if err := c.communicator.Upload(remotePath, input); err != nil {
		errMsg := fmt.Errorf("ssh.UploadData: failed to upload data to remote, err=%v", err)
		res = Result{Cmd: "Upload", Err: errMsg}
		return
	}

	res = Result{Cmd: "Upload"}
	return
}

// UploadFile uploads a local file to the remote path
func (c *Client) UploadFile(remotePath string, file *os.File) (res Result) {
	// ensure we correctly retrieve the output associated with this command
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// print debugging info once the command completes
	defer res.Debug()

	log.Printf("Uploading file (%s) to the remote host at: %s", file.Name(), remotePath)
	if err := c.communicator.Upload(remotePath, file); err != nil {
		errMsg := fmt.Errorf("ssh.UploadFile: failed to upload file to remote host, err=%v", err)
		res = Result{Cmd: "Upload", Err: errMsg}
		return
	}

	res = Result{Cmd: "Upload"}
	return
}

// RunCommand executes a command on the remote host and returns the command's output, the SSH communicator's output and an error
func (c *Client) RunCommand(command string) (res Result) {
	// ensure we correctly retrieve the output associated with this command
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// print debugging info once the command completes
	defer res.Debug()

	// prepare command and output buffer
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := &remote.Cmd{Command: command, Stdout: stdout, Stderr: stderr}

	if err := c.communicator.Start(cmd); err != nil {
		errMsg := fmt.Errorf("ssh.RunCommand:: %v", err)
		res = Result{Cmd: command, Err: errMsg, Stdout: bufferToString(stdout), Stderr: bufferToString(stderr)}
		return
	}

	// await command completion
	if err := cmd.Wait(); err != nil {
		errMsg := fmt.Errorf("cmd.Wait: %v", err)
		res = Result{Cmd: command, Err: errMsg, Stdout: bufferToString(stdout), Stderr: bufferToString(stderr)}
		return
	}

	res = Result{Cmd: command, Stdout: bufferToString(stdout), Stderr: bufferToString(stderr)}
	return
}

// bufferToString read all of a buffer's contents and return a string
func bufferToString(buffer *bytes.Buffer) string {
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		util.PanicOnNonNilErr(fmt.Errorf("ssh.bufferToString: could not save output, err=%v", err))
	}

	return strings.TrimSpace(string(data))
}
