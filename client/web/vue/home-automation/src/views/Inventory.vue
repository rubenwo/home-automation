<template>
    <div class="inventory" style="background-color: rgba(255, 255, 255, 0.60);">
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
            <div slot="actions" slot-scope="props">
                <b-button
                        variant="success"
                        @click="onActionClicked('edit-item', props.rowData)"
                >E
                </b-button
                >
                <b-button
                        variant="danger"
                        @click="onActionClicked('delete-item', props.rowData)"
                >X
                </b-button
                >
            </div>
        </vuetable>
        <div style="padding-top:10px">
            <vuetable-pagination
                    ref="pagination"
                    :css="css.pagination"
                    @vuetable-pagination:change-page="onChangePage"
            ></vuetable-pagination>
        </div>
        <b-button style="width: 100%" @click="addItem()">Add item</b-button>
    </div>
</template>

<script>
  import InventoryService from "../services/inventory.service";
  import Vuetable from "vuetable-2";
  import VuetablePagination from "vuetable-2/src/components/VuetablePagination";
  import _ from "lodash";
  import InventoryFieldsDef from "../services/inventory.fields.def";
  import VuetableBootstrap4Config from "../VuetableBootstrap4Config";
  import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
  // eslint-disable-next-line no-unused-vars
  import fontawesome from "@fortawesome/vue-fontawesome";

  export default {
    name: "Inventory",
    components: {
      Vuetable,
      VuetablePagination,
      // eslint-disable-next-line vue/no-unused-components
      FontAwesomeIcon
    },
    data() {
      return {
        inventory: [],
        addItemData: {
          category: "",
          product: "",
          name: "",
          description: "",
          count: 0
        },
        perPage: 5,
        fields: InventoryFieldsDef,
        css: VuetableBootstrap4Config
      };
    },
    watch: {
      // eslint-disable-next-line no-unused-vars
      inventory(newVal, oldVal) {
        this.$refs.vuetable.refresh();
      }
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
          data: _.slice(local, from, to)
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
          count: 0
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
          count: 0
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
          count: 0
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
          count: 0
        };
        this.$nextTick(() => {
          this.$refs.edit_inventory_item.hide();
        });
      }
    },
    async mounted() {
      const res = await InventoryService.fetchInventory();
      console.log(res);
      this.inventory = res;
    }
  };
</script>

<style scoped>
    .inventory {
        margin-top: 10px;
        text-align: center;
    }
</style>
