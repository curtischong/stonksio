import Pusher from "pusher-js";

export let initWebsockets = () => {
  // Enable pusher logging - don't include this in production

  const pusher = new Pusher("f710317ee72763936d91", {
    cluster: "us2",
  });

  Pusher.logToConsole = false;

  var channel = pusher.subscribe("post");
  channel.bind("new-post", function (data: any) {
    console.log(JSON.stringify(data));
  });
  channel.bind("pusher:subscription_succeeded", function (members: any) {
    console.log("successfully subscribed!");
  });
};
