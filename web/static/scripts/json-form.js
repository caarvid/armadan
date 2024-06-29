htmx.defineExtension("json-form", {
  onEvent: function (name, evt) {
    if (name === "htmx:configRequest") {
      evt.detail.headers["Content-Type"] = "application/json";
    }
  },

  encodeParameters: function (_xhr, params, elt) {
    if (elt.nodeName !== "FORM") {
      return null;
    }

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

    params.entries().forEach(([name, val]) => {
      const e = elt.querySelector(`input[name="${name}"]`);

      if (!e) {
        return;
      }

      switch (e.type) {
        case "number":
          obj[name] = Number(val);
          break;
        case "checkbox":
          obj[name] = val === "on";
        case "hidden":
          if (e.getAttribute("data-type")) {
            obj[name] = Number(val);
          }
          break;
        default:
          obj[name] = val;
      }
    });

    return JSON.stringify(unflatten(obj));
  },
});
