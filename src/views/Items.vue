<template>
  <div class="d-flex flex-row" style="width: 100%">
    <div
      ref="drawerLeft"
      class="flex-grow-0 flex-shrink-0 verticalbar scrollableArea"
      :style="'width: ' + leftWidth + 'px'"
      style="border-right: 1px solid rgba(0, 0, 0, 0.12); overflow-x: hidden !important;"
    >
      <v-expansion-panels accordion multiple>
        <v-expansion-panel>
          <v-expansion-panel-header>
            <v-subheader class="navigationHeader tr-2">
              <transition name="fade-fast">
                <a v-show="leftExtended" class="pl-4">TYPE</a>
              </transition>
              <v-spacer></v-spacer>
              <v-subheader class="tr-2">
                <v-icon :class="{'navigationMenuIcon': !leftExtended}">mdi-format-list-bulleted-type
                </v-icon>
              </v-subheader>
            </v-subheader>
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-list style="padding: 0 !important" dense>
              <hr class="divider"/>
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
                  <v-list-item-icon v-show="leftExtended">
                    <span v-if="'count' in table">{{table['count']}}</span>
                  </v-list-item-icon>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-expansion-panel-content>
        </v-expansion-panel>
        <v-expansion-panel v-show="showLocation">
          <v-expansion-panel-header>
            <v-subheader class="navigationHeader tr-2">
              <transition name="fade-fast">
                <a v-show="leftExtended" class="pl-4">LOCATION</a>
              </transition>
              <v-spacer></v-spacer>
              <v-subheader class="tr-2">
                <v-icon :class="{'navigationMenuIcon': !leftExtended}">mdi-file-tree</v-icon>
              </v-subheader>
            </v-subheader>
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-list style="padding: 0 !important" dense>
              <v-list-item-group>
                <hr class="divider"/>
                <div v-if="loadingTree">
                  <v-spacer/>
                  <looping-rhombuses-spinner
                    style="margin: 0 auto; padding: 20px"
                    :animation-duration="2500"
                    :rhombus-size="10"
                    color="#EF5350"
                  />
                </div>
                <div v-else>
                  <transition name="fade-slow" mode="out-in">
                    <div v-if="!leftExtended" key="dots">
                      <v-icon class="tr-2"
                              style="margin-left: 16px">
                        mdi-dots-horizontal
                      </v-icon>
                    </div>
                    <div v-else key="tree">
                      <v-treeview
                        class="tr-2"
                        v-show="refreshTree && directories.length > 0 && leftExtended"
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
                    </div>
                  </transition>
                </div>
              </v-list-item-group>
            </v-list>
          </v-expansion-panel-content>
        </v-expansion-panel>
        <v-expansion-panel>
          <v-expansion-panel-header>
            <v-subheader class="navigationHeader tr-2">
              <transition name="fade-fast">
                <a v-show="leftExtended" class="pl-4">LABELS</a>
              </transition>
              <v-spacer></v-spacer>
              <v-subheader class="tr-2">
                <v-icon :class="{'navigationMenuIcon': !leftExtended}">mdi-label-outline</v-icon>
              </v-subheader>
            </v-subheader>
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-list style="padding: 0 !important" dense v-show="refreshLabels">
              <v-list-item-group>
                <hr class="divider"/>
                <div v-if="loadingLabels">
                  <v-spacer/>
                  <looping-rhombuses-spinner
                    style="margin: 0 auto; padding: 20px"
                    :animation-duration="2500"
                    :rhombus-size="10"
                    color="#EF5350"
                  />
                </div>
                <div v-else style="margin-top: 10px; margin-left: 10px">
                  <transition name="fade-slow" mode="out-in">
                    <div v-if="!leftExtended" key="dots">
                      <v-icon class="tr-2"
                              style="margin-left: 6px">
                        mdi-dots-horizontal
                      </v-icon>
                    </div>
                    <div v-else key="labels">
                      <v-chip-group
                        v-model="selectedLabel"
                        column
                        active-class="primary--text"
                      >
                        <v-chip
                          class="ma-2 lighten-3 navigationChip"
                          color="grey"
                          text-color="#424242"
                          v-for="label in removeDuplicates(labels)"
                          :key="label"
                          :value="label"
                          small
                          style="margin: 2px !important"
                          @click="filterLabels(label)"
                        >
                          <v-avatar
                            left
                            class="grey lighten-2"
                            style="margin-left: -12px !important; color: #EF5350"
                          >
                            {{countOccurrences(labels, label)}}
                          </v-avatar>
                          {{label}}
                        </v-chip>
                      </v-chip-group>
                    </div>
                  </transition>
                </div>
              </v-list-item-group>
            </v-list>
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-expansion-panels>
      <v-divider/>
      <div class="d-flex justify-end">
        <v-btn small icon @click="toogleLeft">
          <v-icon class="navigationDrawerIcon"
                  :class="{'rotate180': !leftExtended}">mdi-chevron-left
          </v-icon>
        </v-btn>
      </div>
    </div>
    <div class="flex-grow-1 d-flex flex-row" style="width: 100%; overflow: hidden">
      <div style="overflow-x: hidden; transition: width 0.2s ease-in; min-height: 88vh;"
           class="pt-3 flex-shrink-1 flex-grow-1 scrollableArea"
      >
        <div v-if="loadingTable">
          <looping-rhombuses-spinner
            style="margin: 0 auto"
            :animation-duration="2500"
            :rhombus-size="10"
            color="#EF5350"
          />
        </div>
        <div v-else>
          <v-bottom-sheet v-model="bottomSheet"
                          inset
                          hide-overlay
                          no-click-animation
                          persistent>
            <v-sheet class="text-center bottomSheetStyle"
                     height="auto"
                     style="background: #EF5350">
              <v-btn
                small
                icon
                @click="labelsDialogMultiple = true"
              >
                <v-icon class="bottomSheetIcon">mdi-label-outline</v-icon>
              </v-btn>
            </v-sheet>
          </v-bottom-sheet>

          <p>{{selectedLabel}}</p>
          <p>{{filteredLabel}}</p>

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
            v-show="refreshTable"
            v-model="selectedItems"
            :headers="$store.state.headers"
            :items="$store.state.items"
            :loading="loadingFilteredTable"
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
                <td @click="emptyFilter">
                  <v-icon v-if="_.isEmpty(filter.columns) && filteredLabel === ''" color="primary" small class="ml-1">
                    mdi-filter-outline
                  </v-icon>
                  <v-icon v-else color="primary" small class="ml-1">
                    mdi-filter-remove-outline
                  </v-icon>
                </td>
                <td v-for="h in $store.state.headers"
                    :key="h.text" role="columnheader"
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
      </div>
      <div
        :style="'width: ' + rightWidth + '%'"
        style="border-left: 1px solid rgba(0, 0, 0, 0.12)"
        class="flex-grow-0 flex-shrink-0 verticalbar scrollableArea"
        ref="drawerRight"
      >
        <v-dialog
          v-model="labelsDialogMultiple"
          max-width="800px"
          height="500px"
          style="overflow: hidden"
        >
          <v-card>
            <v-card-title class="labelsDialogTitle">
              LABELS
            </v-card-title>
            <hr class="divider" style="margin: 0 !important"/>
            <v-combobox
              v-model="labelsToAddMultiple"
              :items="labels"
              :search-input.sync="searchLabelsMultiple"
              hide-selected
              label="Enter new labels or select existing ones."
              multiple
              chips
              outlined
              style="margin: 20px"
              dense
              deletable-chips
            >
              <template v-slot:no-data>
                <v-list-item>
                  <v-list-item-content>
                    <v-list-item-title>
                      No results matching. Press <kbd>enter</kbd> to
                      create a new one.
                    </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-combobox>
            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="discardAddLabels(true)"
              >
                Close
              </v-btn>
              <v-spacer></v-spacer>
              <v-btn
                color="primary"
                text
                @click="labelsConfirmationMultiple = !labelsConfirmationMultiple"
              >
                Confirm
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <v-dialog
          v-model="labelsConfirmationMultiple"
          max-width="500px"
        >
          <v-card>
            <v-card-title class="labelsDialogTitle">
              <span>Are You Sure?</span>
              <v-spacer></v-spacer>
            </v-card-title>
            <hr class="divider" style="margin: 0 !important"/>
            <div style="padding: 12px">
              <p style="margin: 4px" class="noLabelsText">Following labels will be added to
                {{this.selectedItems.length}}
                items:</p>
              <v-chip style="margin: 4px" v-for="c in labelsToAddMultiple">{{c}}</v-chip>
            </div>
            <hr class="divider" style="margin: 0 !important"/>
            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="labelsConfirmationMultiple = false"
              >
                No
              </v-btn>
              <v-spacer/>
              <v-btn
                color="primary"
                text
                @click="setLabelsMultipleChunk(selectedItems, labelsToAddMultiple)"
              >
                Yes
              </v-btn>
            </v-card-actions>
            <v-progress-linear
              v-if="loadingLabelsOperation"
              color="#EF5350"
              buffer-value="0"
              stream
            ></v-progress-linear>
          </v-card>
        </v-dialog>
        <v-dialog
          v-model="labelsDialog"
          max-width="800px"
          height="500px"
        >
          <v-card>
            <v-card-title class="labelsDialogTitle">
              LABELS
            </v-card-title>
            <hr class="divider" style="margin: 0 !important"/>

            <div
              style="padding: 2px 14px !important; max-height: 40px !important; min-height: 40px !important; background: whitesmoke">
              <v-chip-group
                v-model="labelsToRemove"
                multiple
                show-arrows
              >
                <v-chip class="lighten-3"
                        color="grey"
                        text-color="#EF5350"
                        small
                        filter
                        filter-icon="mdi-close"
                        style="margin: 2px !important"
                        v-for="label in existingLabels" :key="label">
                  <!--                  <v-avatar-->
                  <!--                    left-->
                  <!--                    class="grey lighten-2"-->
                  <!--                    style="margin-left: -33px !important; color: #424242"-->
                  <!--                  >-->
                  <!--                    <v-icon small color="#EF5350">mdi-close</v-icon>-->
                  <!--                  </v-avatar>-->
                  {{ label }}
                </v-chip>
                <p class="noLabelsText"
                   v-if="existingLabels.length === 0 && labelsToAdd.length === 0">No labels
                  exist.</p>
                <v-chip class="lighten-4"
                        color="green"
                        text-color="#EF5350"
                        small
                        filter
                        filter-icon="mdi-close"
                        style="margin: 2px !important"
                        v-for="label in labelsToAdd" :key="label">
                  <!--                  <v-avatar-->
                  <!--                    left-->
                  <!--                    class="grey lighten-2"-->
                  <!--                    style="margin-left: -33px !important; color: #424242"-->
                  <!--                  >-->
                  <!--                    <v-icon small color="#EF5350">mdi-close</v-icon>-->
                  <!--                  </v-avatar>-->
                  {{ label }}
                </v-chip>
              </v-chip-group>
            </div>
            <hr class="divider" style="margin: 0 !important"/>
            <v-combobox
              v-model="labelsToAdd"
              :items="removeExistingLabels(existingLabels)"
              :search-input.sync="searchLabels"
              hide-selected
              label="Enter new labels or select existing ones."
              multiple
              chips
              outlined
              style="margin: 20px"
              dense
              deletable-chips
            >
              <template v-slot:no-data>
                <v-list-item>
                  <v-list-item-content>
                    <v-list-item-title>
                      No results matching. Press <kbd>enter</kbd> to
                      create a new one.
                    </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-combobox>
            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="discardAddLabels(false)"
              >
                Close
              </v-btn>
              <v-spacer></v-spacer>
              <v-btn
                color="primary"
                text
                @click="labelsConfirmation = !labelsConfirmation"
              >
                Confirm
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <v-dialog
          v-model="labelsConfirmation"
          max-width="500px"
        >
          <v-card v-if="this.labelsToRemoveElements">
            <v-card-title class="labelsDialogTitle">
              <span>Are You Sure?</span>
              <v-spacer></v-spacer>
            </v-card-title>

            <div v-if="labelsToAdd.length !== 0">
              <hr class="divider" style="margin: 0 !important"/>
              <div style="padding: 12px">
                <p style="margin: 4px" class="noLabelsText">Following labels will be added:</p>
                <v-chip style="margin: 4px" v-for="c in labelsToAdd">{{c}}</v-chip>
              </div>
            </div>

            <div v-if="calculateLabelsToRemove.length !== 0">
              <hr class="divider" style="margin: 0 !important"/>
              <div style="padding: 12px">
                <p style="margin: 4px" class="noLabelsText">Following labels will be removed:</p>
                <v-chip style="margin: 4px" v-for="c in calculateLabelsToRemove">{{c}}</v-chip>
              </div>
            </div>

            <div v-if="labelsToAdd.length === 0 && calculateLabelsToRemove.length === 0">
              <hr class="divider" style="margin: 0 !important"/>
              <div style="padding: 12px">
                <p style="margin: 4px" class="noLabelsText">No changes made.</p>
              </div>
            </div>

            <hr class="divider" style="margin: 0 !important"/>
            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="labelsConfirmation = false"
              >
                No
              </v-btn>
              <v-spacer/>
              <v-btn
                color="primary"
                text
                @click="setLabelsMultiple($store.state.item.id, labelsToAdd, calculateLabelsToRemove)"
              >
                Yes
              </v-btn>
            </v-card-actions>
            <v-progress-linear
              v-if="loadingLabelsOperation"
              color="#EF5350"
              buffer-value="0"
              stream
            ></v-progress-linear>
          </v-card>
        </v-dialog>
        <v-toolbar
          v-if="!_.isEmpty($store.state.item)"
          class="elevation-0 ml-2 mr-0"
          dense>
          <v-toolbar-title>{{ $store.state.item.id }}</v-toolbar-title>
          <v-spacer/>
          <v-toolbar-items>
            <v-btn
              small
              icon
              @click="labelsDialog = true"
            >
              <v-icon class="detailsIcon">mdi-label-outline</v-icon>
            </v-btn>
            <v-btn
              small
              icon
              @click="expandDetails">
              <v-icon v-if="this.rightWidth === 50" class="detailsIcon">mdi-fullscreen</v-icon>
              <v-icon v-if="this.rightWidth === 100" class="detailsIcon">mdi-fullscreen-exit
              </v-icon>
            </v-btn>
            <v-btn
              small
              icon
              @click="rightNull"
            >
              <v-icon class="detailsIcon">mdi-close</v-icon>
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
  import {LoopingRhombusesSpinner} from 'epic-spinners'
  import {invoke} from "../store/invoke";

  const pause = ms => new Promise(resolve => setTimeout(resolve, ms));

  export default {

    name: 'items',

    components: {
      LoopingRhombusesSpinner,
      item,
    },

    data() {
      return {

        initialLoad: true,

        leftExtended: true,
        itemscol: {},
        rightWidth: 0,
        leftWidth: 250,

        loadingTree: false,
        loadingTable: false,
        loadingLabels: false,
        loadingFilteredTable: false,
        loadingLabelsOperation: false,

        labelsDialog: false,
        labelsConfirmation: false,
        labelsToAdd: [],
        labelsToRemove: [],
        labelsToRemoveElements: [],
        searchLabels: '',

        labelsDialogMultiple: false,
        labelsConfirmationMultiple: false,
        labelsToAddMultiple: [],
        searchLabelsMultiple: '',

        selectedItems: [],
        label: '',
        itemID: '',

        selectedLabel: null,
        filteredLabel: '',

        filter: {
          table: this.$store.state.type,
          columns: {},
        },

        search: '',
        active: [],
        open: [],

        refreshTree: true,
        refreshLabels: true,
        refreshTable: true,

        itemTypeIndex: 0,
        loading: false,
        options: {},

        directories: [],
        labels: [],
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

      showLocation() {
        const type = this.$store.state.type;
        return (type === 'file' || type === 'directory' || type === 'windows-registry-key')
      },

      calculateLabelsToRemove() {

        if (this.$store.state.item.label) {
          const elements = this.existingLabels.concat(this.labelsToAdd);
          const toRemove = [];
          for (let i = 0; i < this.labelsToRemove.length; i++) {
            toRemove.push(elements[this.labelsToRemove[i]]);
          }
          return toRemove
        } else {
          const elements = this.labelsToAdd;
          const toRemove = [];
          for (let i = 0; i < this.labelsToRemove.length; i++) {
            toRemove.push(elements[this.labelsToRemove[i]]);
          }
          return toRemove
        }

      },

      existingLabels() {
        if (this.$store.state.item.label) {
          return Object.keys(this.removeFalseLabels(this.$store.state.item.label));
        } else {
          return [];
        }
      },

      bottomSheet() {
        return this.selectedItems.length !== 0;
      },

    },

    methods: {

      toogleLeft() {
        if (this.leftWidth !== 56) {
          this.leftWidth = 56;
          this.leftExtended = false;
        } else {
          this.leftWidth = 250;
          this.leftExtended = true;
        }
      },

      rightNull() {
        this.rightWidth = 0;
      },
      rightHalf() {
        this.rightWidth = 50;
      },
      rightFull() {
        this.rightWidth = 100;
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

      forceRefreshTree() {
        this.refreshTree = false;
        this.$nextTick(() => {
          this.refreshTree = true;
        });
      },

      forceRefreshLabels() {
        this.refreshLabels = false;
        this.$nextTick(() => {
          this.refreshLabels = true;
        });
      },

      forceRefreshTable() {
        this.refreshTable = false;
        this.$nextTick(() => {
          this.refreshTable = true;
        });
      },

      forceRefreshDetails() {
        this.$store.state.refreshDetails = false;
        console.log("-> ",this.$store.state.refreshDetails)
        this.$nextTick(() => {
          this.$store.state.refreshDetails = true;
          console.log("-> ",this.$store.state.refreshDetails)
        });
      },

      labelRefresher() {

        this.filterLabels(this.filteredLabel);
        this.select(this.$store.state.item);

        this.forceRefreshDetails();
        this.forceRefreshLabels();
        this.forceRefreshTable();

      },

      title(element) {
        if (_.has(this.$store.state.templates, element['type'])) {
          return element[this.$store.state.templates[element['type']].headers[0].value];
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
        return '';
      },

      async getDirectories() {
        this.loadingTree = true;
        const that = this;
        if (this.directories.length === 0) {
          this.$store.dispatch('loadDirectories', {path: ''})
            .then((response) => {
              for (let i = 0; i < response.length; i += 1) {
                if (response[i].name !== '/') {
                  that.directories.push(response[i]);
                }
              }
              this.loadingTree = false;
            })
            .catch(error => console.warn(error));
        }
      },

      async loadList(table) {

        this.loadingTable = true;
        this.loadingTree = true;
        this.initialLoad = true;

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

            this.loadingTable = false;

            this.getDirectories();
            this.getLabels();
            this.forceRefreshTree();

            this.filter = {
              table: this.$store.state.type,
              columns: {},
            };

          });
      },

      emptyFilter() {
        console.log('empty');
        this.loadingFilteredTable = true;
        this.itemscol = {}
        this.filter = {
          table: this.$store.state.type,
          columns: {},
        };
        this.filteredLabel = '';
        this.$store.commit('setLabelFilter', '');
        this.searchFilter()
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
        if (!this.initialLoad) {
          this.loadingFilteredTable = true;
        }
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
        this.$store.dispatch('loadItems').then(() => {
          if (this.initialLoad) {
            this.initialLoad = false;
          } else {
            this.loadingFilteredTable = false;
          }
        })
      },

      updatedir(e) {
        let slash = '';
        let column = '';
        const {type} = this.$store.state;
        this.loadingFilteredTable = true;

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
        this.$store.dispatch('loadItems').then(() => {
          this.loadingFilteredTable = false;
        });
      },

      searchFilter() {
        this.loadingFilteredTable = true;
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

        if (this.search !== '') {
          this.filter.columns.elements = [this.search];
        }
        this.$store.commit('setFilter', this.filter);
        this.$store.dispatch('loadItems').then(() => {
          this.loadingFilteredTable = false;
        });
      },

      clearFilter() {
        this.search = '';
        this.searchFilter();
      },

      expandDetails() {
        if (this.rightWidth === 100) {
          this.rightHalf();
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

      countOccurrences(arr, val) {
        return arr.reduce((a, v) => (v === val ? a + 1 : a), 0);
      },

      removeDuplicates(arr) {
        const filtered = (arr) => arr.filter((v, i) => arr.indexOf(v) === i)
        return filtered(arr);
      },

      removeExistingLabels(arr) {
        return this.labels.filter((elements) => !arr.includes(elements));
      },

      async getLabels(refresh=false) {
        await pause(1500)
        const that = this;
        this.labels = [];
        invoke('GET', '/labels', [], (data) => {
          for (let i = 0; i < data.length; i += 1) {
            that.labels.push(data[i]);
          }
        }).then(() => {
          this.loadingLabels = false;
          if(refresh) {
            this.labelRefresher();
          }
        })
      },

      setLabel(label, id) {
        let url = '/label?id=' + id + '&label=' + label;
        console.log(url)
        invoke('GET', url, [], (data) => {
          console.log(data);
        }).then(() => {
          this.$store.state.item.label[label] = true;
        });
      },

      unsetLabel(label, id) {
        return new Promise((resolve) => {
          let url = '/label?id=' + id + '&label=' + label + '&set=false';
          console.log(url)
          invoke('GET', url, [], (data) => {
            console.log(data);
          }).then(() => {
            this.$store.state.item.label[label] = false;
          });
        });
      },

      async setLabelsMultiple(id, add, remove, chunk=false) {
        if(!chunk) {
          this.loadingLabelsOperation = true;
          console.log("LOADING")
        }
        for (let i = 0; i < add.length; i += 1) {
          this.setLabel(add[i], id)
          await pause(200)
        }
        for (let i = 0; i < remove.length; i += 1) {
          this.unsetLabel(remove[i], id)
          await pause(200)
        }
        if(!chunk) {
          console.log("LOADING FINISHED")
          this.loadingLabelsOperation = false;
          await pause(1000)
          this.labelsConfirmation = false;
          this.getLabels(true);
        }
      },

      async setLabelsMultipleChunk(items, labels) {
        this.loadingLabelsOperation = true;
        for (let i = 0; i < items.length; i += 1) {
          this.setLabelsMultiple(items[i].id, labels, [], true)
        }
        this.loadingLabelsOperation = false;
        await pause(1000)
        this.labelsConfirmationMultiple = false;
        this.getLabels(true);
      },

      filterLabels(label) {
        console.log(label);

        if(label === this.filteredLabel) {
          label = ''
        }

        this.loadingFilteredTable = true;
        this.$store.commit('setLabelFilter', label);
        this.$store.dispatch('loadItems').then(() => {
          this.loadingFilteredTable = false;
          this.filteredLabel = label;
        });
      },

      discardAddLabels(multiple) {
        if (multiple) {
          this.labelsDialogMultiple = false;
          this.labelsToAddMultiple = [];
          this.searchLabelsMultiple = '';
        } else {
          this.labelsDialog = false
          this.labelsToAdd = [];
          this.labelsToRemove = [];
          this.labelsToRemoveElements = [];
          this.searchLabels = '';
        }
      },

      removeFalseLabels(labels) {
        const filtered = Object.keys(labels).reduce(function (filtered, key) {
          if (labels[key] === true) {
            filtered[key] = labels[key];
          }
          return filtered;
        }, {});
        return filtered;
      },

    },

    created() {
      this.loadingTree = true;
      this.loadingTable = true;
      this.loadingLabels = true;
      invoke('GET', '/tables', [], (tables) => {
        this.$store.commit('setTables', tables);
      });
      this.$store.dispatch('loadItems').then(() => {
        this.loadingTree = false;
        this.loadingTable = false;
        this.loadingLabels = false;
      })
    },

    mounted() {
      this.getDirectories();
      this.getLabels();
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

  // *
  //  border-radius: 0 !important

  .verticalbar
    transition: width .2s

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
  // overflow-wrap: anywhere

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

  .navigationDrawerIcon
    transition: $transition-fast
    -moz-transition: $transition-fast
    -webkit-transition: $transition-fast

    &:hover
      background: none !important
      box-shadow: none !important
      transition: $transition-fast
      color: $c-pink !important

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

  .detailsIcon
    padding: 8px

    &:hover
      color: $c-pink !important
      animation: swing
      animation-duration: 0.4s

  .bottomSheetIcon
    padding: 12px
    color: white !important

    &:hover
      color: $c-shadow !important
      animation: swing
      animation-duration: 0.4s

  .content-left
    transition: $transition-fast

  .content-right
    transition: $transition-fast

  .divider
    margin-top: 3px
    margin-bottom: 0
    border: 0
    width: 100%
    border-top: 1px solid $c-pink

  .navigationHeader
    font-size: 14px !important
    font-weight: 400 !important
    padding: 0 !important

  .navigationMenuIcon
    color: $c-pink !important
    transition: $transition-fast !important

  .navigationChip
    &:hover
      cursor: pointer

  .noLabelsText
    padding: 0
    margin: auto
    font-size: 18px
    color: $c-pink
    font-family: 'Roboto Condensed', sans-serif

  .labelsDialogTitle
    font-family: 'Roboto Condensed', sans-serif
    letter-spacing: 0.05rem !important
    font-size: 1.5rem !important
    padding: 16px !important

  .bottomSheetStyle
    border-radius: 10px 10px 0px 0px !important

  .fade-fast-enter,
  .fade-fast-leave-to
    visibility: hidden
    width: 0
    opacity: 0
    padding-left: 0
    overflow-wrap: unset !important
    white-space: nowrap !important

  .fade-fast-enter-active,
  .fade-fast-leave-active
    transition: all 100ms

  .fade-slow-enter,
  .fade-slow-leave-to
    visibility: hidden
    opacity: 0

  .fade-slow-enter-active,
  .fade-slow-leave-active
    transition: all .25s ease-in-out

  .tf-1
    transition: $transition-fastest

  .tf-2
    transition: $transition-faster

  .tf-3
    transition: $transition-fast

  .tf-4
    transition: $transition-slow

  .tf-5
    transition: $transition-slower

  .tf-6
    transition: $transition-slowest

  .o-0
    opacity: 0
    width: 0 !important
    white-space: nowrap !important

  .text-right
    direction: rtl

  .scrollableArea
    overflow: auto !important
    max-height: 94.5vh !important

  .v-expansion-panel-header
    padding: 0 !important
    margin: 0 !important

  .v-expansion-panel-content
    padding: 0 !important
    margin: 0 !important

  .v-expansion-panel-content__wrap
    padding: 0 0 16px !important

  .v-slide-group__next, .v-slide-group__prev
    -ms-flex: 0 1 28px !important
    flex: 0 1 28px !important
    min-width: 28px !important

  .v-chip.v-size--default
    font-size: 12px !important
    height: 24px !important
    padding: 8px !important

  .v-expansion-panel::before
    box-shadow: none !important

  .v-bottom-sheet.v-dialog.v-bottom-sheet--inset
    max-width: fit-content

</style>
