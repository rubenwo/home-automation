(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-522f409a"],{"83a7":function(e,t,n){"use strict";n("d40a")},"8eb4":function(e,t,n){"use strict";n.r(t);var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"routine"},[n("h3",[e._v(e._s(e.routine.name))]),n("p",[e._v("Is active: "+e._s(e.routine.is_active?"True":"False"))]),n("p",[e._v("Trigger: type: "+e._s(e.routine.trigger.type)+", cron_expr: "+e._s(e.routine.trigger.cron_expr))]),n("ul",e._l(e.routine.actions,(function(t,r){return n("li",{key:"action-"+r},[e._v(" ["+e._s(t.method)+"] - ["+e._s(t.addr)+"] "),n("br"),e._v(" Script: ["+e._s(t.script)+"] ")])})),0),n("div",[n("h4",[e._v("Logs")]),null!==e.logs?n("ul",e._l(e.logs,(function(t,r){return n("li",{key:"logs-"+r},[e._v(" ["+e._s(t.logged_at)+"] - ["+e._s(t.message)+"] ")])})),0):e._e()])])},o=[],i=n("1da1"),s=(n("96cf"),n("a1d2")),a={name:"Routine",data:function(){return{routine:{id:0,name:"",is_active:!1,trigger:{cron_expr:"",type:0},actions:[{script:"",data:null,method:"",addr:""}]},logs:[{logged_at:"",message:""}]}},created:function(){var e=this;return Object(i["a"])(regeneratorRuntime.mark((function t(){var n,r;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return e.id=e.$route.params.id,t.next=3,s["a"].fetchRoutine(e.id);case 3:return n=t.sent,console.log(n),e.routine=n.routine,console.log(e.routine),t.next=9,s["a"].fetchLogsForId(e.id);case 9:r=t.sent,console.log(r),e.logs=r.logs,console.log(e.logs);case 13:case"end":return t.stop()}}),t)})))()},mounted:function(){}},u=a,c=(n("83a7"),n("2877")),l=Object(c["a"])(u,r,o,!1,null,"7b8a43ea",null);t["default"]=l.exports},d40a:function(e,t,n){}}]);
//# sourceMappingURL=chunk-522f409a.722f4d51.js.map