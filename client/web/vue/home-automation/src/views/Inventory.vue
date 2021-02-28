<template>
    <div class="inventory" style="background-color: rgba(255, 255, 255, 0.7);">
        <b-modal size="xl"
                 id="add_inventory_item"
                 ref="add_inventory_item"
                 @ok="handleOk"
                 @cancel="handleCancel"
                 @close="handleCancel">
            <b-input-group>
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
                                    v-model="addItemData.count"
                            />
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
        </b-modal>

        <vuetable ref="vuetable"
                  :api-mode="false"
                  :data="inventory"
                  :fields="fields"
                  :data-manager="dataManager"
                  :editable="true"
                  :pagination="true"
                  pagination-path="pagination"
        >
            <div slot="actions" slot-scope="props">
                <!--                <b-button variant="danger" @click="onActionClicked('view-item', props.rowData)">X</b-button>-->
                <b-button variant="success" @click="onActionClicked('edit-item', props.rowData)">E</b-button>
                <b-button variant="danger" @click="onActionClicked('delete-item', props.rowData)">X</b-button>

            </div>
        </vuetable>
        <b-button style="width: 100%" @click="addItem()">Add item</b-button>
    </div>

</template>

<script>
  import InventoryService from '../services/inventory.service';
  import Vuetable from 'vuetable-2'
  import _ from "lodash";

  export default {
    name: "Inventory",
    components: {
      Vuetable
    },
    data() {
      return {
        inventory: [],
        addItemData: {
          name: "",
          description: "",
          count: 0
        },
        fields: [
          {
            name: 'id',
            sortField: 'id'
          },
          {
            name: 'name',
            sortField: 'name'
          },
          {
            name: 'description',
            sortField: 'description'
          },
          {
            name: 'count',
            sortField: 'count'
          },
          "actions"
        ]
      }
    },
    methods: {
      async onActionClicked(event, data) {
        console.log(event)
        console.log(data)
        switch (event) {
          case "edit-item":
            await this.updateItem(data)
            break
          case "delete-item":
            await this.removeItem(data.id)
            break
        }
      },
      async dataManager(sortOrder, pagination) {
        if (this.data.length < 1) return;

        let local = this.data;

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
        console.log('pagination:', pagination)
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
        const id = data.id;
        const res = await InventoryService.updateInventoryItem(id, data);
        console.log(res);
        this.inventory = res;
      },
      async handleOk(evt) {
        evt.preventDefault();
        console.log(this.addItemData);
        this.handleSubmit();
      },
      async handleCancel() {
        console.log("CANCEL")
      },
      async handleSubmit() {
        const res = await InventoryService.addInventoryItem(this.addItemData);
        console.log(res);
        this.inventory = res;
        this.addItemData = {
          name: "",
          description: "",
          count: 0
        };
        this.$nextTick(() => {
          this.$refs.deviceModal.hide();
        });
      }
    },
    async mounted() {
      const res = await InventoryService.fetchInventory();
      console.log(res)
      this.inventory = res;
    },
  }
</script>

<style scoped>
    .inventory {
        margin-top: 10px;
        text-align: center;
    }
</style>
