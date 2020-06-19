<template>
  <div>
    <v-data-table
      :headers="headers"
      :items="logs"
      :fixed-header="true"
      :footer-props="{'items-per-page-options': [50, 100]}"
      :items-per-page="10"
      style="overflow: visible !important;"
      show-select
      dense
    >
      <template v-slot:body.prepend>
        <tr>
          <td>
            <v-icon
              v-if="(timeFilter === '' ) && (fileFilter === '' ) && (messageFilter === '' )"
              color="primary"
              small class="ml-1">
              mdi-filter-outline
            </v-icon>
            <v-icon v-else color="primary" small class="ml-1" @click="emptyFilter">
              mdi-filter-remove-outline
            </v-icon>
          </td>
          <td v-for="h in headers"
              :key="h.text" role="columnheader"
              scope="col">
            <v-text-field
              v-model="itemscol[h.value]"
              @keyup.enter.native="searchFilter"
              hide-details
              label="Filter"
              clearable
            />
          </td>
        </tr>
      </template>
      <template v-slot:item.time="{ item }">
        <div>{{ toLocal(item.time) }}</div>
      </template>
    </v-data-table>
  </div>
</template>

<script>
  import {invoke} from "../store/invoke";
  import {DateTime} from "luxon";

  export default {
    name: 'logs',
    data() {
      return {

        timeFilter: '',
        fileFilter: '',
        messageFilter: '',

        headers: [
          {
            text: "Time",
            value: "time",
            filter: f => {
              if (this.timeFilter === '' || this.timeFilter === null) return true
              else return (f + '').toLowerCase().includes(this['timeFilter'].toLowerCase())
            }
          },
          {
            text: "File",
            value: "file",
            filter: f => {
              if (this.fileFilter === '' || this.fileFilter === null) return true
              else return (f + '').toLowerCase().includes(this['fileFilter'].toLowerCase())
            }
          },
          {
            text: "Message",
            value: "message",
            filter: f => {
              if (this.messageFilter === '' || this.messageFilter === null) return true
              else return (f + '').toLowerCase().includes(this['messageFilter'].toLowerCase())
            }
          },
        ],
        logs: [],
        itemscol: {},
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

      emptyFilter() {
        console.log('empty');
        Vue.$set(this, 'itemscol', {});
        this.searchFilter();
      },

      searchFilter() {

        let column = '';

        for (const key in this.itemscol) {
          const value = this.itemscol[key];
          column = key;

          if (column === 'time') {
            this.timeFilter = value
          }

          if (column === 'file') {
            this.fileFilter = value
          }

          if (column === 'message') {
            this.messageFilter = value
          }
        }

      },

    },

    mounted() {
      this.loadFiles();
    },
  };
</script>
