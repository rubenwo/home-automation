// Bootstrap4 + FontAwesome Icon
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

export default {
  components: {
    FontAwesomeIcon,
  },
  table: {
    tableWrapper: "",
    tableHeaderClass: "mb-0",
    tableBodyClass: "mb-0",
    tableClass: "table table-bordered table-hover",
    loadingClass: "loading",
    ascendingIcon: "fa fa-chevron-up",
    descendingIcon: "fa fa-chevron-down",
    ascendingClass: "sorted-asc",
    descendingClass: "sorted-desc",
    sortableIcon: "fa fa-sort",
    detailRowClass: "vuetable-detail-row",
    handleIcon: "fa fa-bars text-secondary",
    // eslint-disable-next-line no-unused-vars
    renderIcon(classes, options) {
      return `<i class="${classes.join(" ")}"></span>`;
    },
  },
  pagination: {
    wrapperClass: "pagination float-right",
    activeClass: "active",
    disabledClass: "disabled",
    pageClass: "page-item",
    linkClass: "page-link",
    paginationClass: "pagination",
    paginationInfoClass: "float-left",
    dropdownClass: "form-control",
    icons: {
      first: "fa fa-chevron-left",
      prev: "fa fa-chevron-left",
      next: "fa fa-chevron-right",
      last: "fa fa-chevron-right",
    },
  },
};
