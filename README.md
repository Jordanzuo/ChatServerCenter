# ChatServerCenter
ChatServerCenter、ChatServer、ChatServerModel、ChatClient是一个聊天组合；共同提供聊天服务。
其中：
ChatServerCenter：聊天的中心服务器，用于服务器调度，消息转发，对外提供API等功能。是聊天服务器体系的核心。
ChatServer：聊天服务器，为客户端提供连接、登陆、发送消息等功能；是与客户端直接连接的服务器。
ChatServerModel：定义ChatServerCenter和ChatServer公用的对象。
ChatClient：测试使用。
此系统支持动态扩展，即可以根据客户端的数量来扩展ChatServer的数量。但是ChatServerCenter是不可扩展的。

备注：
本系统与ChatServer_Go, ChatClient_Go是不同的。他们是一个提供聊天服务的组合，但是不支持动态扩展。
