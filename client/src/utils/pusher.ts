
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
  pusherInstance = new Pusher('8dbf2e2ddc742f692e39', {
    cluster: 'us2'
  });
  return pusherInstance;
};

export default pusher;
