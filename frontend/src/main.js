import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import './assets/main.css'; // Your global plain CSS

const app = createApp(App);

app.use(createPinia()); // Initialize Pinia first so Router can use it
app.use(router);

app.mount('#app');