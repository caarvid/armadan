htmx.defineExtension("json-form", {
  onEvent: function(name, evt) {
    if (name === "htmx:configRequest") {
      evt.detail.headers["Content-Type"] = "application/json";
    }
  },

  encodeParameters: function(_xhr, params, elt) {
    if (elt.nodeName !== "FORM") {
      return null;
    }

    function unflatten(data) {
      var result = {};
      for (var i in data) {
        var keys = i.split(".");

        keys.reduce(function(r, e, j) {
          return (
            r[e] ||
            (r[e] = isNaN(Number(keys[j + 1]))
              ? keys.length - 1 == j
                ? data[i]
                : {}
              : [])
          );
        }, result);
      }

      return result;
    }

    for (var name of Object.keys(params)) {
      const e = elt.querySelector(`input[name="${name}"]`);

      if (!e) {
        continue;
      }

      switch (e.type) {
        case "number":
          params[name] = Number(params[name]);
          break;
        case "hidden":
          if (e.getAttribute("data-type")) {
            params[name] = Number(params[name]);
          }
          break;
      }
    }

    return JSON.stringify(unflatten(params));
  },
});
