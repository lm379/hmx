import { defineComponent } from 'vue';
import NavBar from '../components/NavBar.vue';
import ChatWindow from '../components/ChatWindow.vue';

export default defineComponent({
  name: 'App',
  components: {
    NavBar,
    ChatWindow
  }
});
