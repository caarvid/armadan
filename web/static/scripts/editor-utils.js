/**
 * @param {HTMLTextAreaElement} el
 * @param {(t: string) => string} textFn
 * @param {{ focusEmpty: boolean }} config
 */
function textAction(el, textFn, { focusEmpty } = { focusEmpty: true }) {
  const text = el.value.substring(el.selectionStart, el.selectionEnd);
  const newText = textFn(text);

  el.setRangeText(newText, el.selectionStart, el.selectionEnd, "select");

  if (focusEmpty && !text.length) {
    const middle =
      el.selectionStart + (el.selectionEnd - el.selectionStart) / 2;
    el.setSelectionRange(middle, middle);
  }

  el.focus();
}

/**
 * @param {string} s
 * @returns {(t: string) => string}
 */
function surround(s) {
  return (text) => {
    if (text.startsWith(s) && text.endsWith(s)) {
      return text.substring(s.length, text.length - s.length);
    }

    return s + text + s;
  };
}
