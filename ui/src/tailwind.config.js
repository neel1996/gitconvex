// tailwind.config.js
module.exports = {
  purge: {
    content: ["./src/**/*.html", "./src/**/*.js"],
    options: {
      whitelist: [
        "hover:bg-blue-500",
        "hover:bg-green-500",
        "hover:bg-pink-500",
        "border-yellow-900",
        "bg-yellow-200",
        "border-red-900",
        "bg-red-200",
        "border-green-900",
        "bg-green-200",
      ],
    },
  },
  theme: {},
  variants: {},
  plugins: [],
};
