import React from "react";
import { useSpring, animated } from "react-spring";

export default function InfiniteLoader({ loadAnimation }) {
  let infiniteLoader = useSpring({
    from: { marginLeft: "0px", opacity: 0 },
    to: async (next) => {
      let i = 100;
      while (--i) {
        await next({
          marginLeft: "10px",
          opacity: 0.5,
        });

        await next({
          marginLeft: "20px",
          opacity: 0.8,
        });

        await next({
          marginLeft: "30px",
          opacity: 1,
        });

        await next({
          marginLeft: "0px",
          opacity: 1,
        });
      }
    },
    config: {
      tension: "500",
    },
  });
  return (
    <div className="flex gap-4 mx-auto text-center">
      {loadAnimation ? (
        <>
          <animated.div
            className="bg-pink-200 w-8 h-8 rounded-full p-2"
            style={infiniteLoader}
          ></animated.div>
          <animated.div
            className="bg-green-200 w-8 h-8 rounded-full p-2"
            style={infiniteLoader}
          ></animated.div>
          <animated.div
            className="bg-blue-200 w-8 h-8 rounded-full p-2"
            style={infiniteLoader}
          ></animated.div>
        </>
      ) : null}
    </div>
  );
}
