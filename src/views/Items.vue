<template>
  <div>
    <v-navigation-drawer
      ref="drawerLeft"
      color="sidebar"
      dark
      clipped
      app
      :style="'flex-basis: '+navigationLeft.width+'px'"
      :width="navigationLeft.width"
      v-model="navigationLeft.shown"
      style="overflow: hidden"
    >
      <div class="mt-2">
        <v-list dense>
          <v-subheader>Type</v-subheader>
          <v-list-item-group v-model="itemTypeIndex" color="primary">
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
                <v-icon v-else>mdi-help</v-icon>
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
      </div>
      <div class="mt-2">
        <v-list dense>
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
      <v-col
        align-self="end"
      >
          <v-icon class="navigationDrawerIcon"
                  :class="{'rotate180': navigationLeft.mini}"
                  @click="toogleLeft">mdi-chevron-left
          </v-icon>
      </v-col>

    </v-navigation-drawer>
    <v-container fluid :class="{'pr-0': !this.detailsExpanded && this.detailsHidden}">
      <v-row class="ma-0" style="flex-wrap: nowrap;">
        <!--        <v-col style="display: -webkit-box !important; width: 100%; flex-basis: 100% !important"-->
        <!--               class="flex-grow-1 flex-shrink-1 pb-3 content-left pa-0"-->
        <!--               v-show="!detailsExpanded">-->
        <v-col class="flex-shrink-1 pb-3 content-left pa-0"
               :class="dataTable">
          <v-card class="pt-3">
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
              class="mx-3 my-0 py-0"
            ></v-text-field>
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
              show-select
              style="overflow: visible !important;"
              dense
            >
              <template v-slot:body.prepend>
                <tr>
                  <td>
                    <v-icon style="font-size: 16px; margin-left: 4px; color: #ED3C54">
                      mdi-filter-outline
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
          </v-card>
        </v-col>
        <v-col
          :style="'overflow: hidden; width: ' + navigationRight.width + 'px; flex-basis: ' + navigationRight.width + 'px'"
          ref="drawerRight"
          class="ma-0 pt-0 pr-0 content-right"
          transition="scale-transition"
          :class="{ 'pl-0': detailsExpanded }"
        >
          <v-card>
            <v-toolbar
              v-if="!_.isEmpty($store.state.item) && !navigationRight.mini"
              class="elevation-0 ml-2 mr-0"
              dense>
              <v-toolbar-title>Details</v-toolbar-title>
              <v-spacer></v-spacer>
              <v-toolbar-items>
                <v-btn
                  small
                  icon
                  @click="expandDetails">
                  <v-icon class="detailsIcon">mdi-arrow-all</v-icon>
                </v-btn>
                <v-btn
                  small
                  icon
                  @click="closeLeft"
                >
                  <v-icon class="detailsIcon">mdi-close</v-icon>
                </v-btn>
              </v-toolbar-items>
            </v-toolbar>
            <v-divider class="mx-0"></v-divider>
            <transition name="scale-transition">
              <item v-show="!navigationRight.mini" :content="$store.state.item"/>
            </transition>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
  import { DateTime } from 'luxon';
  import item from '@/views/Document.vue';
  import Vuetify from '../plugins/vuetify';

  const pause = ms => new Promise(resolve => setTimeout(resolve, ms));

  export default {

    name: 'items',

    components: {
      item,
    },

    data() {
      return {

        itemscol: {},
        detailsExpanded: false,
        detailsHidden: true,

        filter: {
          table: this.$store.state.type,
          columns: {},
        },

        search: '',
        active: [],
        open: [],
        refresh: true,

        navigationLeft: {
          type: 'left',
          mini: false,
          shown: true,
          width: 250,
          oldWidth: 250,
          borderSize: 3,
        },

        navigationRight: {
          type: 'right',
          mini: true,
          shown: true,
          width: 0,
          oldWidth: 500,
          borderSize: 3,
        },

        itemTypeIndex: 0,
        loading: false,
        options: {},

        directories: [],
      };
    },

    computed: {

      dataTable() {
        return {
          'dataTableVisible': !this.detailsExpanded,
          'flex-grow-1': !this.detailsExpanded,
          'dataTableInvisible': this.detailsExpanded
        };
      },

      theme() {
        return (Vuetify.framework.theme.dark) ? 'dark' : 'light';
      },

      dark() {
        return Vuetify.framework.theme.dark;
      },

      light() {
        return !Vuetify.framework.theme.dark;
      },

      count() {
        for (let i = 0; i < this.$store.state.tables.length; i += 1) {
          if (this.$store.state.type === this.$store.state.tables[i].name) {
            return this.$store.state.tables[i].count;
          }
        }
        return 0;
      },

      paneSize: {
        get() {
          return this.$store.state.listPane;
        },
        set(value) {
          this.$store.commit('setlistPane', value);
        },
      },

    },

    methods: {

      toggle(navigation) {
        if (!navigation.mini) {
          navigation.oldWidth = parseInt(navigation.width, 10);
          if (navigation.type === 'left') {
            navigation.width = 56;
          } else {
            navigation.width = 0;
          }
        }
        navigation.mini = !navigation.mini;

        if (!navigation.mini) {
          navigation.width = parseInt(navigation.oldWidth, 10);
        }
      },

      toogleLeft() {
        this.toggle(this.navigationLeft);
      },

      toogleRight() {
        this.detailsHidden = !this.detailsHidden;
        this.toggle(this.navigationRight);
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

      closeLeft() {
        if (this.detailsExpanded) {
          this.toogleRight();
          this.detailsExpanded = !this.detailsExpanded;
        } else {
          this.toogleRight();
        }
      },

      async getDirectories() {
        const that = this;
        if (this.directories.length === 0) {
          this.$store.dispatch('loadDirectories', { path: '' })
            .then((response) => {
              for (let i = 0; i < response.length; i += 1) {
                if (response[i].name !== '/') {
                  that.directories.push(response[i]);
                }
              }
            })
            .catch(error => console.warn(error));

          console.log('INITIAL: ', JSON.stringify(this.directories, null, 1));
        }
      },

      async loadList(table) {

        this.detailsExpanded = false;
        this.detailsHidden = true;
        this.directories = [];
        this.itemscol = {};
        this.active = [];
        this.open = [];

        this.$store.commit('setItem', {});
        if (!this.navigationRight.mini) {
          this.toogleRight();
        }
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

      select(e) {
        this.$store.commit('setItem', e);
        if (this.navigationRight.mini) {
          this.toogleRight();
        }
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
        const { type } = this.$store.state;

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
        if (!this.detailsExpanded) {
          this.detailsExpanded = true;
        } else {
          this.detailsExpanded = false;
        }
      },

      async fetchTreeChildren(item) {
        const directories = [];

        this.$store.dispatch('loadDirectories', { path: item.path })
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
  @import '~animate.css'

  table tr td
    cursor: pointer !important

  *
    border-radius: 0 !important

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
    overflow-wrap: anywhere

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

  .navigationDrawerIcon
    transition: $transition-slower
    -moz-transition: 0.5s ease-in-out
    -webkit-transition: 0.5s ease-in-out
    display: inline-block
    position: absolute
    bottom: 0
    right: 0
    width: 100%
    background: none !important
    box-shadow: none
    font-size: 30px
    margin: 0 !important

    &:hover
      color: $c-pink

  .rotate180
    -ms-transform-origin: 50% 50%
    -webkit-transform-origin: 50% 50%
    -moz-transform-origin: 50% 50%
    transform-origin: 50% 50%
    transform: rotate(180deg)
    -moz-transform: rotate(180deg)
    -webkit-transform: rotate(180deg)
    -o-transform: rotate(180deg)
    -ms-transform: rotate(180deg)

  /*.navigationDrawerIconHoverNarrow
    &:hover
      color: $c-pink !important
      animation: swing !important
      animation-duration: 0.5s

  /*.navigationDrawerIconHoverWide
    &:hover
      color: $c-pink !important
      animation: swing !important
      animation-duration: 0.5s

  .dataTableVisible
    width: 100%
    flex-basis: 100% !important

  .dataTableInvisible
    width: 0 !important
    flex-basis: 0 !important
    flex-grow: 0 !important

  .detailsIcon
    &:hover
      color: $c-pink !important
      animation: swing
      animation-duration: 0.4s

</style>
