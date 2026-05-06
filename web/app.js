const endpoints = {
  validate: "/api/validate",
  minify: "/api/minify",
  escape: "/api/escape",
  unescape: "/api/unescape",
  unicodeDecode: "/api/unicode/decode",
  unicodeEncode: "/api/unicode/encode",
};

const editor = document.querySelector("#editor");
const statusNode = document.querySelector("#status");

function setStatus(message, isError = false) {
  statusNode.textContent = message;
  statusNode.style.color = isError ? "#b42318" : "#5b6b7a";
}

async function runAction(action) {
  const endpoint = endpoints[action];
  if (!endpoint) {
    return;
  }

  setStatus("处理中...");
  try {
    const response = await fetch(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ input: editor.value }),
    });
    const payload = await response.json();
    if (!response.ok || !payload.ok) {
      throw new Error(payload.error || "处理失败");
    }
    editor.value = payload.result;
    setStatus("完成");
  } catch (error) {
    setStatus(error.message || "处理失败", true);
  }
}

document.querySelectorAll("[data-action]").forEach((button) => {
  button.addEventListener("click", () => runAction(button.dataset.action));
});

document.querySelector("#copy").addEventListener("click", async () => {
  try {
    await navigator.clipboard.writeText(editor.value);
    setStatus("结果已复制");
  } catch (error) {
    setStatus("复制失败", true);
  }
});

document.querySelector("#clear").addEventListener("click", () => {
  editor.value = "";
  setStatus("已清空");
});

document.querySelector("#saveFile").addEventListener("click", () => {
  const blob = new Blob([editor.value], { type: "application/json;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = "ejson-result.json";
  anchor.click();
  URL.revokeObjectURL(url);
  setStatus("已开始下载");
});
