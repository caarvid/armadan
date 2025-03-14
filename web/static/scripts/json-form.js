htmx.defineExtension("json-form", {
  onEvent: function (name, evt) {
    if (name === "htmx:configRequest") {
      evt.detail.headers["Content-Type"] = "application/json";
    }
  },

  encodeParameters: function (_xhr, params, _elt) {
    function unflatten(data) {
      const result = {};

      Object.entries(data).forEach(([name, val]) => {
        const keys = name.split(".");

        keys.reduce(function (r, e, j) {
          return (
            r[e] ||
            (r[e] = isNaN(Number(keys[j + 1]))
              ? keys.length - 1 == j
                ? val
                : {}
              : [])
          );
        }, result);
      });

      return result;
    }

    const obj = {};

    for (const [name, val] of params.entries()) {
      const e = document.body.querySelector(`[name="${name}"]`);

      if (!e) {
        return;
      }

      switch (e.type) {
        case "number":
          obj[name] = Number(val);
          break;
        case "checkbox":
          obj[name] = val === "on";
          break;
        case "hidden":
          if (e.getAttribute("data-type")) {
            obj[name] = Number(val);
          } else {
            obj[name] = val;
          }
          break;
        default:
          obj[name] = val;
      }
    }

    return JSON.stringify(unflatten(obj));
  },
});
