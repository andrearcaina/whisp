function chatApp() {
    return {
        messages: [],
        newMessage: "",
        connected: false,
        ws: null,

        showGifModal: false,
        selectedGif: null,
        gifs: [],
        gifSearchQuery: "",
        gifLoading: false,
        searchHistory: JSON.parse(
            localStorage.getItem("gifSearchHistory") || "[]",
        ),
        currentGifCategory: "trending",
        searchTimeout: null,
        gifCategories: [
            { id: "trending", name: "Trending", icon: "🔥" },
            { id: "reaction", name: "Reactions", icon: "😀" },
            { id: "sports", name: "Sports", icon: "⚽" },
            { id: "stickers", name: "Stickers", icon: "✨" },
        ],

        init() {
            this.connectWebSocket();
            this.loadMessages();
            this.loadTrendingGifs();
        },

        connectWebSocket() {
            // check if ws is production or development
            const isProduction = window.location.hostname !== "localhost";
            const wsProtocol = isProduction ? "wss" : "ws";
            const wsHost = isProduction
                ? window.location.hostname
                : "localhost:8080";
            this.ws = new WebSocket(`${wsProtocol}://${wsHost}/ws`);

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
            if (
                (!this.newMessage.trim() && !this.selectedGif) ||
                !this.connected
            )
                return;
            const message = {
                message: this.newMessage.trim(),
                username: "anonymous",
                gif_url: this.selectedGif?.gif_url || null,
            };
            this.ws.send(JSON.stringify(message));
            this.newMessage = "";
            this.selectedGif = null;
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
            } catch {
                this.messages = [];
            }
        },

        scrollToBottom() {
            this.$nextTick(() => {
                this.$refs.messagesContainer.scrollTop =
                    this.$refs.messagesContainer.scrollHeight;
            });
        },

        formatTime(timestamp) {
            return timestamp ? new Date(timestamp).toLocaleString() : "";
        },

        openGifModal() {
            this.showGifModal = true;
            if (!this.gifs.length && this.currentGifCategory === "trending")
                this.loadTrendingGifs();
        },

        closeGifModal() {
            this.showGifModal = false;
            this.gifSearchQuery = "";
            if (this.searchTimeout) {
                clearTimeout(this.searchTimeout);
                this.searchTimeout = null;
            }
        },

        onSearchInput() {
            if (this.searchTimeout) clearTimeout(this.searchTimeout);
            if (!this.gifSearchQuery.trim()) {
                this.currentGifCategory = "trending";
                this.loadTrendingGifs();
                return;
            }
            this.searchTimeout = setTimeout(() => this.performSearch(), 500);
        },

        onSearchEnter() {
            if (this.searchTimeout) {
                clearTimeout(this.searchTimeout);
                this.searchTimeout = null;
            }
            if (this.gifSearchQuery.trim()) this.performSearch();
        },

        async performSearch() {
            if (!this.gifSearchQuery.trim()) return;
            this.currentGifCategory = "search";
            this.gifLoading = true;
            try {
                const response = await fetch(
                    `/api/klipy/gifs/${encodeURIComponent(this.gifSearchQuery)}`,
                );
                if (response.ok) {
                    const data = await response.json();
                    this.gifs = data.map((gif) => ({
                        id: gif.id,
                        gif_url: gif.gif_url,
                        title: gif.title || this.gifSearchQuery,
                    }));
                    this.addToSearchHistory(this.gifSearchQuery);
                } else {
                    this.gifs = [];
                }
            } catch {
                this.gifs = [];
            }
            this.gifLoading = false;
        },

        async loadTrendingGifs() {
            this.gifLoading = true;
            try {
                const response = await fetch("/api/klipy/gifs/trending");
                if (response.ok) {
                    const data = await response.json();
                    this.gifs = data.map((gif) => ({
                        id: gif.id,
                        gif_url: gif.gif_url,
                        title: gif.title || "Trending GIF",
                    }));
                } else {
                    this.gifs = [];
                }
            } catch {
                this.gifs = [];
            }
            this.gifLoading = false;
        },

        async switchGifCategory(categoryId) {
            this.currentGifCategory = categoryId;
            this.gifSearchQuery = "";
            this.gifLoading = true;
            try {
                let searchTerm = "";
                if (categoryId === "trending") {
                    await this.loadTrendingGifs();
                    return;
                } else if (categoryId === "reaction") searchTerm = "reaction";
                else if (categoryId === "sports") searchTerm = "sports";
                else if (categoryId === "stickers") searchTerm = "stickers";

                const response = await fetch(`/api/klipy/gifs/${searchTerm}`);
                if (response.ok) {
                    const data = await response.json();
                    this.gifs = data.map((gif) => ({
                        id: gif.id,
                        gif_url: gif.gif_url,
                        title: gif.title || searchTerm,
                    }));
                } else {
                    this.gifs = [];
                }
            } catch {
                this.gifs = [];
            }
            this.gifLoading = false;
        },

        selectGif(gif) {
            this.selectedGif = gif;
            this.closeGifModal();
        },

        removeSelectedGif() {
            this.selectedGif = null;
        },

        addToSearchHistory(query) {
            if (!query.trim()) return;
            this.searchHistory = this.searchHistory.filter(
                (item) => item !== query,
            );
            this.searchHistory.unshift(query);
            this.searchHistory = this.searchHistory.slice(0, 10);
            localStorage.setItem(
                "gifSearchHistory",
                JSON.stringify(this.searchHistory),
            );
        },

        useSearchHistory(query) {
            this.gifSearchQuery = query;
            this.performSearch();
        },
    };
}
