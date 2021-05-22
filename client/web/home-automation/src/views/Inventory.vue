<template>
  <div class="inventory" style="background-color: rgba(255, 255, 255, 0.6)">
    <b-modal
      size="xl"
      id="add_inventory_item"
      ref="add_inventory_item"
      @ok="handleOk"
      @cancel="handleCancel"
      @close="handleCancel"
    >
      <b-input-group>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Category:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="Category"
                type="text"
                v-model="addItemData.category"
              />
            </b-col>
          </b-row>
        </b-container>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Product:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="Product"
                type="text"
                v-model="addItemData.product"
              />
            </b-col>
          </b-row>
        </b-container>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Name:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="name"
                type="text"
                v-model="addItemData.name"
              />
            </b-col>
          </b-row>
        </b-container>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Description:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="description"
                v-model="addItemData.description"
              />
            </b-col>
          </b-row>
        </b-container>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Count:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="name"
                type="number"
                v-model.number="addItemData.count"
              />
            </b-col>
          </b-row>
        </b-container>
      </b-input-group>
    </b-modal>
    <b-modal
      size="xl"
      id="edit_inventory_item"
      ref="edit_inventory_item"
      @ok="handleEditOk"
      @cancel="handleEditCancel"
      @close="handleEditCancel"
    >
      <b-input-group>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Category:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="Category"
                type="text"
                v-model="addItemData.category"
              />
            </b-col>
          </b-row>
        </b-container>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Product:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="Product"
                type="text"
                v-model="addItemData.product"
              />
            </b-col>
          </b-row>
        </b-container>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Name:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="name"
                type="text"
                v-model="addItemData.name"
              />
            </b-col>
          </b-row>
        </b-container>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Description:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="description"
                v-model="addItemData.description"
              />
            </b-col>
          </b-row>
        </b-container>
        <b-container fluid>
          <b-row class="my-1">
            <b-col sm="4">
              <label>Count:</label>
            </b-col>
            <b-col sm="8">
              <b-form-input
                size="sm"
                class="mx-1"
                placeholder="name"
                type="number"
                v-model.number="addItemData.count"
              />
            </b-col>
          </b-row>
        </b-container>
      </b-input-group>
    </b-modal>
    <filter-bar />

    <vuetable
      ref="vuetable"
      :api-mode="false"
      :fields="fields"
      :per-page="perPage"
      :css="css.table"
      :data-manager="dataManager"
      pagination-path="pagination"
      @vuetable:pagination-data="onPaginationData"
    >
      <div slot="actions" scope="props">
        <div class="table-button-container">
          <b-button
                  class="btn btn-warning btn-sm"
                  @click="onActionClicked('edit-item', props.rowData)"
          >Edit
          </b-button>
          <b-button
                  class="btn btn-danger btn-sm"
                  @click="onActionClicked('delete-item', props.rowData)"
          >Delete
          </b-button>
        </div>
      </div>

    </vuetable>
    <div style="padding-top: 10px">
      <vuetable-pagination-info
              ref="paginationInfoTop"
      />

      <vuetable-pagination
        ref="pagination"
        :css="css.pagination"
        @vuetable-pagination:change-page="onChangePage"
      />
    </div>
    <b-button style="width: 100%" @click="addItem()">Add item</b-button>
  </div>
</template>

<script>
import InventoryService from "../services/inventory.service";
import Vuetable from "vuetable-2";
import VuetablePagination from "vuetable-2/src/components/VuetablePagination";
import VuetablePaginationInfo from "vuetable-2/src/components/VuetablePaginationInfo";
import _ from "lodash";
import InventoryFieldsDef from "../services/inventory.fields.def";
// import VuetableBootstrap4Config from "../components/VuetableBootstrap4Config";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import FilterBar from "../components/FilterBar";

export default {
  name: "Inventory",
  components: {
    Vuetable,
    VuetablePagination,
    // eslint-disable-next-line vue/no-unused-components
    FontAwesomeIcon,
    FilterBar,
    VuetablePaginationInfo,
  },
  data() {
    return {
      inventory: [],
      fullInventory: [],
      addItemData: {
        category: "",
        product: "",
        name: "",
        description: "",
        count: 0,
      },
      perPage: 5,
      fields: InventoryFieldsDef,
      css: {
        table: {
          tableClass: 'table table-striped table-bordered table-hovered',
          loadingClass: 'loading',
          ascendingIcon: 'glyphicon glyphicon-chevron-up',
          descendingIcon: 'glyphicon glyphicon-chevron-down',
          handleIcon: 'glyphicon glyphicon-menu-hamburger',
        },
        pagination: {
          infoClass: 'pull-left',
          wrapperClass: 'vuetable-pagination pull-right',
          activeClass: 'btn-primary',
          disabledClass: 'disabled',
          pageClass: 'btn btn-border',
          linkClass: 'btn btn-border',
          icons: {
            first: '',
            prev: '',
            next: '',
            last: '',
          },
        }}
    };
  },
  watch: {
    // eslint-disable-next-line no-unused-vars
    inventory(newVal, oldVal) {
      this.$refs.vuetable.refresh();
    },
  },
  methods: {
    onPaginationData(paginationData) {
      this.$refs.pagination.setPaginationData(paginationData);
    },
    onChangePage(page) {
      this.$refs.vuetable.changePage(page);
    },

    async onActionClicked(event, data) {
      console.log(event);
      console.log(data);
      switch (event) {
        case "edit-item":
          await this.updateItem(data);
          break;
        case "delete-item":
          await this.removeItem(data.id);
          break;
      }
    },
    dataManager(sortOrder, pagination) {
      if (this.inventory.length < 1) return;

      let local = this.inventory;

      // sortOrder can be empty, so we have to check for that as well
      if (sortOrder.length > 0) {
        console.log("orderBy:", sortOrder[0].sortField, sortOrder[0].direction);
        local = _.orderBy(
          local,
          sortOrder[0].sortField,
          sortOrder[0].direction
        );
      }

      pagination = this.$refs.vuetable.makePagination(
        local.length,
        this.perPage
      );
      console.log("pagination:", pagination);
      let from = pagination.from - 1;
      let to = from + this.perPage;

      return {
        pagination: pagination,
        data: _.slice(local, from, to),
      };
    },
    async removeItem(id) {
      console.log("removing item with id: " + id);
      const res = await InventoryService.deleteInventoryItem(id);
      console.log(res);
      this.inventory = res;
    },
    async addItem() {
      this.$refs.add_inventory_item.show();
    },
    async updateItem(data) {
      this.addItemData = data;
      this.$refs.edit_inventory_item.show();
    },
    async handleOk(evt) {
      evt.preventDefault();
      console.log(this.addItemData);
      this.handleSubmit();
    },
    async handleCancel() {
      this.addItemData = {
        name: "",
        description: "",
        count: 0,
      };
    },
    async handleSubmit() {
      console.log(this.addItemData);
      const res = await InventoryService.addInventoryItem(this.addItemData);
      console.log(res);
      this.inventory = res;
      this.addItemData = {
        name: "",
        description: "",
        count: 0,
      };
      this.$nextTick(() => {
        this.$refs.add_inventory_item.hide();
      });
    },
    async handleEditOk(evt) {
      evt.preventDefault();
      console.log(this.addItemData);
      this.handleEditSubmit();
    },
    async handleEditCancel() {
      this.addItemData = {
        name: "",
        description: "",
        count: 0,
      };
    },
    async handleEditSubmit() {
      console.log(this.addItemData);
      const res = await InventoryService.updateInventoryItem(
        this.addItemData.id,
        this.addItemData
      );
      console.log(res);
      this.inventory = res;
      this.addItemData = {
        name: "",
        description: "",
        count: 0,
      };
      this.$nextTick(() => {
        this.$refs.edit_inventory_item.hide();
      });
    },
    async onFilterSet(filterText) {
      const res = await InventoryService.fetchInventory();
      const filtered = res.filter((item) => {
        if (item.category.toLowerCase().includes(filterText)) return true;
        if (item.product.toLowerCase().includes(filterText)) return true;
        if (item.name.toLowerCase().includes(filterText)) return true;
        if (item.description.toLowerCase().includes(filterText)) return true;
        return false;
      });
      if (filtered.length === 0) {
        console.log("EMPTY");
        this.inventory = [];
        this.fields = [];
      } else {
        this.fields = InventoryFieldsDef;
        this.inventory = filtered;
      }
    },
    async onFilterReset() {
      const res = await InventoryService.fetchInventory();
      this.fields = InventoryFieldsDef;
      this.inventory = res;
    },
  },
  async mounted() {
    const res = await InventoryService.fetchInventory();
    console.log(res);
    this.inventory = res;
    this.$events.$on("filter-set", (eventData) => this.onFilterSet(eventData));
    // eslint-disable-next-line no-unused-vars
    this.$events.$on("filter-reset", (e) => this.onFilterReset());
  },
};
</script>

<style scoped>
.inventory {
  margin-top: 10px;
  text-align: center;
}
.orange.glyphicon {
  color: orange;
}

th.sortable {
  color: #ec971f;
}

</style>
