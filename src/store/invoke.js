import axios from 'axios';

export function invoke(method, url, arg, callback) {
  return new Promise( resolve => {
    if (window.external.invoke !== undefined) {
      // local mode
      window.external.invoke(JSON.stringify(arg));
    } else {
      // server mode
      let server = '';
      if (typeof webpackHotUpdate !== 'undefined') {
        // dev mode
        server = 'http://' + window.location.hostname + ':8081';
      }
      axios({
        method: method,
        url: server + '/api' + url,
        data: arg
      })
        .then(response => {
          callback(response.data);
          resolve(response.data);
        });
    }
  });
}
