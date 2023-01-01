import React from "react";
import FontAwesome from "@fortawesome/fontawesome-free";

import "./Chat.css";

export default function Chat({ ...rest }) {
  console.log(FontAwesome);
  return (
    <div {...rest} className={`${rest["className"] || ""} chatComponent`}>
      <div className="chatHistory">
        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut viverra nunc
        sem, a dapibus sem sodales ac. Phasellus laoreet lectus sit amet nibh
        luctus, vel convallis ante molestie. Proin ut quam lacinia, accumsan
        ipsum at, tristique metus. Integer vel scelerisque enim. In gravida
        fermentum ligula eget tempus. Curabitur at nunc ultrices, vestibulum
        enim et, convallis sem. Aenean tellus nulla, venenatis a imperdiet ut,
        molestie eu tellus. Sed mauris sapien, porta in neque in, fringilla
        luctus ipsum. Quisque scelerisque congue dictum. Duis quam nisi, ornare
        sed sem ac, maximus feugiat tortor. Sed sed gravida velit. Quisque
        accumsan fringilla vestibulum. Mauris non ex neque. Mauris sed tempus
        eros. Proin condimentum, felis nec lobortis tristique, neque enim
        facilisis nulla, non sodales ex nisl nec odio. Nulla facilisi.
        Pellentesque non sapien et ex consequat rhoncus. Nullam quis sem luctus
        enim consequat semper. Sed et vestibulum augue. Aenean tempor auctor
        leo, non tincidunt lorem pulvinar vitae. Donec sodales libero eu nisi
        egestas ultricies. Aenean at sem felis. In et leo velit. Aliquam quam
        augue, hendrerit a neque quis, fermentum rhoncus massa. Pellentesque
        habitant morbi tristique senectus et netus et malesuada fames ac turpis
        egestas. Nulla bibendum, massa a congue blandit, quam lacus congue odio,
        vitae semper orci est eu neque. Sed tincidunt eget neque quis rutrum.
        Nunc quis sapien vitae massa congue consequat. Nunc feugiat, ipsum non
        fringilla bibendum, libero quam elementum diam, nec varius odio mauris
        at nisi. In quis dictum ligula, tincidunt fringilla ligula. Cras ut nunc
        elementum, dapibus nunc id, auctor dui. Quisque vel dolor bibendum,
        ullamcorper dolor at, fermentum nisl. Cras malesuada ante in ultrices
        laoreet. Ut eu sem id sapien maximus vulputate. Curabitur pulvinar vel
        magna et imperdiet. Proin id arcu faucibus, tempor tellus placerat,
        finibus lorem. Donec id felis et metus laoreet viverra. Quisque dapibus
        feugiat ullamcorper. In semper vel dui in semper. Sed mattis dolor et
        mauris lacinia, sed faucibus augue dapibus. Curabitur mi nunc, congue id
        dictum ac, pulvinar non lorem. Nulla ut turpis lorem. Nulla semper
        finibus felis, non porttitor metus feugiat vitae. Aliquam erat volutpat.
        Pellentesque eget fringilla mi, in auctor ipsum. Aliquam vel lacus in
        tortor ultricies fermentum. Curabitur ut tellus accumsan ex pulvinar
        iaculis. Fusce eget neque ex. Duis et felis venenatis, condimentum quam
        ac, maximus velit. Proin gravida vitae augue non consectetur. Nulla sit
        amet suscipit libero. Integer posuere elit tempor sodales accumsan. Sed
        luctus, est eu tristique ultricies, ipsum mauris feugiat arcu, dignissim
        vestibulum risus lacus a enim. Mauris tempor ex vitae lorem rhoncus, nec
        ultricies eros mollis. Phasellus vulputate tempor feugiat. Nunc
        fringilla magna et lobortis viverra. Duis dolor orci, ultricies ut
        feugiat malesuada, porttitor non mi. Praesent sit amet quam mollis,
        tempor sem et, tincidunt nunc.
      </div>
      <div className="chatInput">
        <input placeholder="Enter chat message..." type="text"></input>
        <button>
          <i className="fa-solid fa-paper-plane"></i>
        </button>
      </div>
    </div>
  );
}
