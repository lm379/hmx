import { defineComponent, ref, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

export default defineComponent({
  name: 'RegisterView',
  setup() {
    const router = useRouter();
    const username = ref('');
    const phone = ref('');
    const email = ref('');
    const code = ref('');
    const password = ref('');
    const loading = ref(false);
    const sending = ref(false);
    const error = ref('');
    const success = ref('');
    const countdown = ref(0);
    let timer: any = null;

    const startCountdown = () => {
      countdown.value = 60;
      timer = setInterval(() => {
        countdown.value--;
        if (countdown.value <= 0) {
          clearInterval(timer);
        }
      }, 1000);
    };

    const sendCode = async () => {
      if (!email.value) {
        error.value = '请先输入邮箱';
        return;
      }
      try {
        sending.value = true;
        error.value = '';
        await axios.post('/api/v1/auth/send-code', { email: email.value });
        success.value = '验证码已发送，请查收邮箱';
        startCountdown();
      } catch (err: any) {
        error.value = err.response?.data?.message || '发送失败，请重试';
      } finally {
        sending.value = false;
      }
    };

    const handleRegister = async () => {
      try {
        loading.value = true;
        error.value = '';
        success.value = '';
        await axios.post('/api/v1/auth/register', {
          username: username.value,
          phone: phone.value,
          email: email.value,
          code: code.value,
          password: password.value
        });

        success.value = '注册成功！正在跳转登录...';
        setTimeout(() => {
          router.push('/login');
        }, 1500);
      } catch (err: any) {
        error.value = err.response?.data?.message || '注册失败，请检查输入';
      } finally {
        loading.value = false;
      }
    };

    onUnmounted(() => {
      if (timer) clearInterval(timer);
    });

    return {
      username,
      phone,
      email,
      code,
      password,
      loading,
      sending,
      error,
      success,
      countdown,
      sendCode,
      handleRegister
    };
  }
});
