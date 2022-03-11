import axios from "axios";

class API {
  async listResources() {
    const result = await axios.get("/api/resources");
    return result.data;
  }

  async addCluster(data) {
    const result = await axios.post("/api/clusters", data);
    return result.data;
  }
  async editCluster(data) {
    const result = await axios.put("/api/clusters", data);
    return result.data;
  }
  async deleteCluster({ cluster_id }, permanent) {
    const result = await axios.delete(`/api/clusters/${cluster_id}${permanent ? "?permanent=true" : ""}`);
    return result.data;
  }

  async addFaaS(data) {
    const result = await axios.post("/api/faas-deployments", data);
    return result.data;
  }
  async editFaaS(data) {
    const result = await axios.put("/api/faas-deployments", data);
    return result.data;
  }
  async deleteFaaS({ faas_id }) {
    const result = await axios.delete(`/api/faas-deployments/${faas_id}`);
    return result.data;
  }

  async addStorage(data) {
    const result = await axios.post("/api/storage-deployments", data);
    return result.data;
  }
  async deleteStorage({ storage_id }, permanent) {
    const result = await axios.delete(`/api/storage-deployments/${storage_id}${permanent ? "?permanent=true" : ""}`);
    return result.data;
  }

  async addBucket(data) {
    const result = await axios.post("/api/buckets", data);
    return result.data;
  }
  async editBucket(data) {
    const result = await axios.put("/api/buckets", data);
    return result.data;
  }
  async deleteBucket({ bucket_id }) {
    const result = await axios.delete(`/api/buckets/${bucket_id}`);
    return result.data;
  }

  async addObject(data) {
    const result = await axios.post("/api/objects", data);
    return result.data;
  }
  async deleteObject({ object_id }) {
    const result = await axios.delete(`/api/objects/${object_id}`);
    return result.data;
  }

  async setLBDefaults(data) {
    return axios.put("/api/load-balancer/settings", data);
  }
  async setLBOverrides(data) {
    return axios.put("/api/load-balancer/route-overrides", data);
  }
}

const api = new API();

export default api;
