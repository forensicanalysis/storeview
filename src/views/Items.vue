<template>
  <div class="d-flex flex-row" style="width: 100%">
    <div
      ref="drawerLeft"
      class="flex-grow-0 flex-shrink-0 verticalbar"
      :style="'width: ' + leftWidth + 'px'"
      style="border-right: 1px solid rgba(0, 0, 0, 0.12)"
    >
      <div class="d-flex justify-end">
        <v-btn small icon @click="toogleLeft">
          <v-icon>mdi-chevron-left</v-icon>
        </v-btn>
      </div>
      <v-list dense>
        <v-subheader>Type</v-subheader>
        <v-list-item-group v-model="itemTypeIndex" color="primary">
          <v-list-item @click="loadList('')">
            <v-list-item-icon>
              <v-icon>mdi-file-multiple</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>All</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            v-for="(table, i) in _.sortBy($store.state.tables, ['name'])"
            :key="i"
            @click="loadList(table['name'])"
            class="px-4"
          >
            <v-list-item-icon>
              <v-icon
                v-if="_.has($store.state.templates[table['name']], 'icon')"
                v-text="'mdi-'+$store.state.templates[table['name']].icon"/>
              <v-icon v-else>mdi-file-outline</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title v-text="_.startCase(table.name)"/>
            </v-list-item-content>
            <v-list-item-icon>
              <span v-if="'count' in table">{{table['count']}}</span>
            </v-list-item-icon>
          </v-list-item>
        </v-list-item-group>
      </v-list>
      <v-list class="mt-2" dense>
        <v-list-item-group>
          <v-subheader>Location</v-subheader>
          <v-treeview
            v-if="refresh && directories.length > 0"
            activatable
            hoverable
            dense
            transition
            @update:active="updatedir"
            :active.sync="active"
            :items="directories"
            :load-children="fetchTreeChildren"
            :open.sync="open"
            item-key="path"
            color="primary"
            ref="treeView"
          >
            <template v-slot:prepend="{ item, open }">
              <v-icon v-if="item.children">
                {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
              </v-icon>
              <v-icon v-if="!item.children">mdi-folder-outline</v-icon>
            </template>
          </v-treeview>
        </v-list-item-group>
      </v-list>
    </div>
    <div class="flex-grow-1 d-flex flex-row" style="width: 100%; overflow: hidden">
      <div style="overflow-x: scroll; transition: width 0.2s ease-in"
           class="pt-3 flex-shrink-1 flex-grow-1"
      >
        <v-text-field
          v-model="search"
          prepend-inner-icon="mdi-magnify"
          clear-icon="mdi-close-circle"
          clearable
          dense
          outlined
          @keyup.enter.native="searchFilter"
          @click:append="searchFilter"
          @click:clear="clearFilter"
          class="mx-3"
        />
        <v-data-table
          :headers="$store.state.headers"
          :items="$store.state.items"
          :loading="loading"
          :options.sync="options"
          :server-items-length="$store.state.itemCount"
          @update:options="updateopt"
          :fixed-header="true"
          @click:row="select"
          :footer-props="{'items-per-page-options': [10, 25, 50, 100]}"
          :items-per-page="50"
          show-select
          style="overflow: visible !important;"
          dense
        >
          <template v-slot:body.prepend>
            <tr>
              <td>
                <v-icon v-if="_.isEmpty(filter.columns)" color="primary" small class="ml-1">
                  mdi-filter-outline
                </v-icon>
                <v-icon v-else color="primary" small class="ml-1" @click="emptyFilter">
                  mdi-filter-remove-outline
                </v-icon>
              </td>
              <td v-for="h in $store.state.headers" :key="h.text" role="columnheader"
                  scope="col">
                <v-text-field
                  v-model="itemscol[h.value]"
                  @keyup.enter.native="searchFilter"
                  hide-details
                  label="Filter"
                  :reverse="h.align === 'right'"
                  clearable
                />
              </td>
            </tr>
          </template>

          <template v-slot:item.title="{ item }">
            <div>{{ title(item) }}</div>
          </template>
          <template v-slot:item.atime="{ item }">
            <div>{{ toLocal(item.atime) }}</div>
          </template>
          <template v-slot:item.mtime="{ item }">
            <div>{{ toLocal(item.mtime) }}</div>
          </template>
          <template v-slot:item.ctime="{ item }">
            <div>{{ toLocal(item.ctime) }}</div>
          </template>
          <template v-slot:item.origin="{ item }">
            <div><span v-if="'path' in item.origin">{{item.origin.path}}</span></div>
          </template>
          <template v-slot:item.size="{ item }">
            <div>{{ humanBytes(item.size, true) }}</div>
          </template>
        </v-data-table>
      </div>
      <div
        :style="'width: ' + rightWidth + '%'"
        style="border-left: 1px solid rgba(0, 0, 0, 0.12)"
        class="flex-grow-0 flex-shrink-0 verticalbar"
        ref="drawerRight"
      >
        <v-toolbar
          v-if="!_.isEmpty($store.state.item)"
          class="elevation-0 mx-2"
          dense>
          <v-toolbar-title>{{ $store.state.item.id }}</v-toolbar-title>
          <v-spacer/>
          <v-toolbar-items>
            <v-btn
              small
              icon
              @click="expandDetails">
              <v-icon>mdi-fullscreen</v-icon>
            </v-btn>
            <v-btn
              small
              icon
              @click="rightNull"
            >
              <v-icon>mdi-close</v-icon>
            </v-btn>
          </v-toolbar-items>
        </v-toolbar>
        <v-divider class="mx-0"/>
        <item :content="$store.state.item"/>
      </div>
    </div>
  </div>
</template>

<script>
  import {DateTime} from 'luxon';
  import item from '@/views/Document.vue';

  const pause = ms => new Promise(resolve => setTimeout(resolve, ms));

  export default {

    name: 'items',

    components: {
      item,
    },

    data() {
      return {

        itemscol: {},
        rightWidth: 0,
        leftWidth: 250,

        filter: {
          table: this.$store.state.type,
          columns: {},
        },

        search: '',
        active: [],
        open: [],
        refresh: true,

        itemTypeIndex: 0,
        loading: false,
        options: {},

        directories: [],
      };
    },

    computed: {

      count() {
        for (let i = 0; i < this.$store.state.tables.length; i += 1) {
          if (this.$store.state.type === this.$store.state.tables[i].name) {
            return this.$store.state.tables[i].count;
          }
        }
        return 0;
      },

    },

    methods: {

      toogleLeft() {
        if (this.leftWidth !== 56) {
          this.leftWidth = 56
        } else {
          this.leftWidth = 250
        }
      },

      rightNull() {
        this.rightWidth = 0
      },
      rightHalf() {
        this.rightWidth = 50
      },
      rightFull() {
        this.rightWidth = 100
      },

      humanBytes(bytes, si) {
        const thresh = si ? 1000 : 1024;
        if (Math.abs(bytes) < thresh) {
          return `${bytes}B`;
        }
        const units = si
          ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
          : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
        let u = -1;
        do {
          bytes /= thresh;
          u += 1;
        } while (Math.abs(bytes) >= thresh && u < units.length - 1);
        return `${bytes.toFixed(1)}${units[u]}`;
      },

      forceRefresh() {
        this.refresh = false;
        this.$nextTick(() => {
          this.refresh = true;
        });
      },

      title(element) {
        if (_.has(this.$store.state.templates, element['type'])) {
          return element[this.$store.state.templates[element['type']].headers[0].value]
        }
        if (_.has(element, 'name')) {
          return element['name'];
        }
        if (_.has(element, 'key')) {
          return element['key'];
        }
        if (_.has(element, 'title')) {
          return element['title'];
        }
        return "";
      },

      async getDirectories() {
        const that = this;
        if (this.directories.length === 0) {
          this.$store.dispatch('loadDirectories', {path: ''})
            .then((response) => {
              for (let i = 0; i < response.length; i += 1) {
                if (response[i].name !== '/') {
                  that.directories.push(response[i]);
                }
              }
            })
            .catch(error => console.warn(error));
        }
      },

      async loadList(table) {

        this.directories = [];
        this.itemscol = {};
        this.active = [];
        this.open = [];

        this.$store.commit('setItem', {});
        this.rightNull();
        this.$store.commit('setTable', table);
        this.$store.dispatch('loadItems')
          .then(() => {
            for (let i = 0; i < this.$store.state.headers.length; i += 1) {
              this.itemscol[this.$store.state.headers[i].value] = '';
            }

            this.getDirectories();
            this.forceRefresh();

            this.filter = {
              table: this.$store.state.type,
              columns: {},
            };
          });
      },

      emptyFilter() {
        console.log("empty");
        Vue.$set(this, "itemscol", {});
        this.searchFilter();
      },

      select(e) {
        this.$store.commit('setItem', e);
        this.rightHalf();
      },

      toLocal(s) {
        return DateTime.fromISO(s)
          .toLocaleString(DateTime.DATETIME_SHORT_WITH_SECONDS);
      },

      updateopt(e) {
        this.$store.commit('setOffset', (e.page - 1) * e.itemsPerPage);
        this.$store.commit('setLimit', e.itemsPerPage);
        const sort = {
          table: this.$store.state.type,
          columns: {},
        };
        for (let i = 0; i < e.sortBy.length; i += 1) {
          if (e.sortDesc[i]) {
            sort.columns[e.sortBy[i]] = 'DESC';
          } else {
            sort.columns[e.sortBy[i]] = 'ASC';
          }
        }

        this.$store.commit('setSort', sort);
        this.$store.dispatch('loadItems');
      },

      updatedir(e) {
        let slash = '';
        let column = '';
        const {type} = this.$store.state;

        if (type === 'directory') {
          slash = '/';
          column = 'path';
        } else if (type === 'file') {
          slash = '/';
          column = 'origin.path';
        } else if (type === 'windows-registry-key') {
          slash = '\\';
          column = 'key';
        } else {
          console.log('TABLE DOES NOT EXIST!');
        }

        if (e.length === 0) {
          this.filter.table = type;
          this.filter.columns[column] = [];
        } else {
          this.filter.table = type;
          this.filter.columns[column] = [e[0] + slash];
        }

        this.$store.commit('setFilter', this.filter);
        this.$store.dispatch('loadItems');
      },

      searchFilter() {
        this.filter.table = this.$store.state.type;

        let column = '';

        for (const key in this.itemscol) {
          const value = this.itemscol[key];
          if (key === 'origin') {
            column = 'origin.path';
          } else {
            column = key;
          }
          this.filter.columns[column] = [value];
        }

        this.filter.columns.elements = [this.search];
        this.$store.commit('setFilter', this.filter);
        this.$store.dispatch('loadItems');
      },

      clearFilter() {
        this.search = '';
        this.searchFilter();
      },

      expandDetails() {
        if (this.rightWidth === 100) {
          this.rightHalf()
        } else {
          this.rightFull();
        }
      },

      async fetchTreeChildren(item) {
        const directories = [];

        this.$store.dispatch('loadDirectories', {path: item.path})
          .then((response) => {
            if ((response.length === 1) && (response[0].name === '/')) {
              delete item.children;
            } else {
              for (let i = 0; i < response.length; i += 1) {
                if (response[i].name !== '/') {
                  directories.push(response[i]);
                }
              }
            }
          })
          .catch(error => console.warn(error));

        await pause(1000);

        const key = item.path;
        const parentNode = this.$refs.treeView.nodes[key];

        let childNode;
        directories.forEach((child) => {
          childNode = {
            ...parentNode,
            item: child,
            vnode: null,
          };
          this.$refs.treeView.nodes[child.path] = childNode;
        });

        for (let i = 0; i < directories.length; i += 1) {
          item.children.push(directories[i]);
        }
      },

    },

    mounted() {
      this.getDirectories();
    },

    destroyed() {
      this.$store.commit('setItem', {});
    },

  };


</script>

<style lang="sass">
  @import '~vuetify/src/styles/styles.sass'
  @import '../styles/colors.scss'
  @import '../styles/animations.scss'

  table tr td
    cursor: pointer !important

  *
    border-radius: 0 !important

  .verticalbar
      transition: width .2s

  .fade-enter-active, .fade-leave-active
    transition: opacity .2s

  .fade-enter, .fade-leave-to
  /* .fade-leave-active below version 2.1.8 */
     opacity: 0

  .v-treeview-node__label, .v-data-table .v-text-field, .v-data-table .v-label
    font-size: 0.8125rem

  .v-data-table .v-label--active
    opacity: 0

  .v-data-table .v-text-field > .v-input__control > .v-input__slot:before
    border-style: none

  .v-text-field
    padding-top: 0
    margin-top: 0

  .v-data-table--fixed-header .v-data-table__wrapper
    overflow-y: hidden

  .v-list-item__icon
    min-width: max-content
    font-size: 12px
    line-height: 2

  .v-application--is-ltr .v-list-item__icon:first-child
    margin-right: 8px
    margin-top: 6px

  .v-list--dense .v-subheader
    font-size: 0.75rem
    font-weight: bold

  .v-input .v-label
    font-size: 12px

  .content-left
    transition: $transition-fast

  .content-right
    transition: $transition-fast

</style>
