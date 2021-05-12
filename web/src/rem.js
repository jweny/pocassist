// 设计稿以1920px为宽度，页面宽度设计为10rem

const baseSize = 192; // 这个是设计稿中1rem的大小。
function setRem() {
  // 实际设备页面宽度和设计稿的比值
  const scale = document.documentElement.clientWidth / 1920;
  // 计算实际的rem值并赋予给html的font-size
  document.documentElement.style.fontSize = baseSize * scale + "px";
}
setRem();
window.addEventListener("resize", () => {
  setRem();
});
