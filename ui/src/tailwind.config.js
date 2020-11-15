// tailwind.config.js
module.exports = {
  purge: {
    content: ["./src/**/*.html", "./src/**/*.js"],
    options: {
      whitelist: ["hover:bg-blue-500", "hover:bg-green-500", "hover:bg-pink-500"],
    },
  },
  theme: {},
  variants: {},
  plugins: [],
};
