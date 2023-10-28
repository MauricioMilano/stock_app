package chatmessage_model

type ChatMessage struct {
	ChatMessage  string `json:"chatMessage"`
	ChatUser     string `json:"chatUser"`
	ChatRoomId   uint   `json:"chatRoomId"`
	ChatRoomName string `json:"chatRoomName"`
}
