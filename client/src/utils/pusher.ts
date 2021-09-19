
import Pusher from 'pusher-js';

const DISCONNECTED = 'disconnected';

let pusherInstance: Pusher;

const pusher = (): Pusher => {
  if (pusherInstance) {
    if (pusherInstance.connection.state === DISCONNECTED) {
      pusherInstance.connect();
    }
    return pusherInstance;
  }
  pusherInstance = new Pusher('f710317ee72763936d91', {
    cluster: 'us2'
  });
  return pusherInstance;
};

export default pusher;
