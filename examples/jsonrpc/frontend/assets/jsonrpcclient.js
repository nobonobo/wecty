const defaultHeaders = {
  Accept: "application/json",
  "Content-Type": "application/json",
};

class JsonRpcClient {
  constructor(endpoint = "/rpc/", headers = {}) {
    this.lastId = 0;
    this.endpoint = endpoint;
    this.headers = Object.assign({}, defaultHeaders, headers);
  }
  request(method, ...params) {
    const id = this.lastId++;
    const req = {
      method: "POST",
      headers: this.headers,
      credential: true,
      body: JSON.stringify({
        jsonrpc: "2.0",
        id,
        method,
        params: Array.isArray(params) ? params : [params],
      }),
    };
    return fetch(this.endpoint, req)
      .then((res) => checkStatus(res))
      .then((res) => parseJSON(res))
      .then((res) => checkError(res, req))
      .then((res) => logResponse(res));
  }
}

function parseJSON(response) {
  return response.json();
}

function checkStatus(response) {
  // we assume 400 as valid code here because it's the default return code when sth has gone wrong,
  // but then we have an error within the response, no?
  if (response.status >= 200 && response.status <= 400) {
    return response;
  }

  const error = new Error(response.statusText);
  error.response = response;
  throw error;
}

function checkError(data, req) {
  if (data.error) {
    const error = new RpcError(data.error, req, data);
    error.response = data;
    throw error;
  }
  return data;
}

function logResponse(response) {
  return response.result;
}

/**
 * RpcError is a simple error wrapper holding the request and the response.
 */
class RpcError extends Error {
  constructor(message, request, response) {
    super(message);

    this.name = "RpcError";
    this.message = message || "";
    this.request = request;
    this.response = response;
  }
  toString() {
    return this.message;
  }
  getRequest() {
    return this.request;
  }
  getResponse() {
    return this.response;
  }
}

export { JsonRpcClient, RpcError };
