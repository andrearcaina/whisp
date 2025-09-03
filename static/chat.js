function chatApp() {
	return {
		messages: [],
		newMessage: "",
		connected: false,
		ws: null,

		init() {
			this.connectWebSocket();
			this.loadMessages();
		},

		connectWebSocket() {
			this.ws = new WebSocket("ws://localhost:8080/ws");

			this.ws.onopen = () => {
				this.connected = true;
			};

			this.ws.onclose = () => {
				this.connected = false;
				setTimeout(() => this.connectWebSocket(), 2000);
			};

			this.ws.onmessage = (event) => {
				const message = JSON.parse(event.data);
				this.messages.push(message);
				this.scrollToBottom();
			};
		},

		sendMessage() {
			if (!this.newMessage.trim() || !this.connected) return;

			const message = {
				message: this.newMessage.trim(),
				username: "anonymous"
			};

			this.ws.send(JSON.stringify(message));
			this.newMessage = "";
		},

		async loadMessages() {
			try {
				const response = await fetch("/api/messages");
				if (response.ok) {
					const data = await response.json();
					this.messages = data || [];
					this.messages.reverse();
					this.scrollToBottom();
				} else {
					this.messages = [];
				}
			} catch (error) {
				this.messages = [];
			}
		},

		scrollToBottom() {
			this.$nextTick(() => {
				this.$refs.messagesContainer.scrollTop = this.$refs.messagesContainer.scrollHeight;
			});
		},

		formatTime(timestamp) {
			if (!timestamp) return '';
			return new Date(timestamp).toLocaleTimeString();
		}
	}
}
