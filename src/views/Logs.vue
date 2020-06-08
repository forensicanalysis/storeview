<template>
  <v-data-table
    :headers="headers"
    :items="logs"
    :fixed-header="true"
    :footer-props="{'items-per-page-options': [50, 100]}"
    :items-per-page="10"
    style="overflow: visible !important;"
    dense
  >
    <template v-slot:item.time="{ item }">
      <div>{{ toLocal(item.time) }}</div>
    </template>
  </v-data-table>
</template>

<script>
  import {invoke} from "../store/invoke";
  import {DateTime} from "luxon";

  export default {
    name: 'logs',
    data() {
      return {
        headers: [
          {text: "Time", value: "time"},
          {text: "File", value: "file"},
          {text: "Message", value: "message"},
        ],
        logs: [],
      };
    },

    methods: {
      loadFiles() {
        invoke('GET', '/logs', [], (data) => {
          console.log(data);
          this.logs = data.elements;
        });
      },

      toLocal(s) {
        return DateTime.fromISO(s)
          .toLocaleString(DateTime.DATETIME_SHORT_WITH_SECONDS);
      },
    },

    mounted() {
      this.loadFiles();
    },
  };
</script>
