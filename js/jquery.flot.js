/* Javascript plotting library for jQuery, version 0.8.3.

Copyright (c) 2007-2014 IOLA and Ole Laursen.
Licensed under the MIT license.

*/

// first an inline dependency, jquery.colorhelpers.js, we inline it here
// for convenience

/* Plugin for jQuery for working with colors.
 *
 * Version 1.1.
 *
 * Inspiration from jQuery color animation plugin by John Resig.
 *
 * Released under the MIT license by Ole Laursen, October 2009.
 *
 * Examples:
 *
 *   $.color.parse("#fff").scale('rgb', 0.25).add('a', -0.5).toString()
 *   var c = $.color.extract($("#mydiv"), 'background-color');
 *   console.log(c.r, c.g, c.b, c.a);
 *   $.color.make(100, 50, 25, 0.4).toString() // returns "rgba(100,50,25,0.4)"
 *
 * Note that .scale() and .add() return the same modified object
 * instead of making a new one.
 *
 * V. 1.1: Fix error handling so e.g. parsing an empty string does
 * produce a color rather than just crashing.
 */
!function(a){a.color={},a.color.make=function(t,e,i,o){var n={};return n.r=t||0,n.g=e||0,n.b=i||0,n.a=null!=o?o:1,n.add=function(t,e){for(var i=0;i<t.length;++i)n[t.charAt(i)]+=e;return n.normalize()},n.scale=function(t,e){for(var i=0;i<t.length;++i)n[t.charAt(i)]*=e;return n.normalize()},n.toString=function(){return 1<=n.a?"rgb("+[n.r,n.g,n.b].join(",")+")":"rgba("+[n.r,n.g,n.b,n.a].join(",")+")"},n.normalize=function(){function t(t,e,i){return e<t?t:i<e?i:e}return n.r=t(0,parseInt(n.r),255),n.g=t(0,parseInt(n.g),255),n.b=t(0,parseInt(n.b),255),n.a=t(0,n.a,1),n},n.clone=function(){return a.color.make(n.r,n.b,n.g,n.a)},n.normalize()},a.color.extract=function(t,e){for(var i;(""==(i=t.css(e).toLowerCase())||"transparent"==i)&&((t=t.parent()).length&&!a.nodeName(t.get(0),"body")););return"rgba(0, 0, 0, 0)"==i&&(i="transparent"),a.color.parse(i)},a.color.parse=function(t){var e,i=a.color.make;if(e=/rgb\(\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*\)/.exec(t))return i(parseInt(e[1],10),parseInt(e[2],10),parseInt(e[3],10));if(e=/rgba\(\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*,\s*([0-9]+(?:\.[0-9]+)?)\s*\)/.exec(t))return i(parseInt(e[1],10),parseInt(e[2],10),parseInt(e[3],10),parseFloat(e[4]));if(e=/rgb\(\s*([0-9]+(?:\.[0-9]+)?)\%\s*,\s*([0-9]+(?:\.[0-9]+)?)\%\s*,\s*([0-9]+(?:\.[0-9]+)?)\%\s*\)/.exec(t))return i(2.55*parseFloat(e[1]),2.55*parseFloat(e[2]),2.55*parseFloat(e[3]));if(e=/rgba\(\s*([0-9]+(?:\.[0-9]+)?)\%\s*,\s*([0-9]+(?:\.[0-9]+)?)\%\s*,\s*([0-9]+(?:\.[0-9]+)?)\%\s*,\s*([0-9]+(?:\.[0-9]+)?)\s*\)/.exec(t))return i(2.55*parseFloat(e[1]),2.55*parseFloat(e[2]),2.55*parseFloat(e[3]),parseFloat(e[4]));if(e=/#([a-fA-F0-9]{2})([a-fA-F0-9]{2})([a-fA-F0-9]{2})/.exec(t))return i(parseInt(e[1],16),parseInt(e[2],16),parseInt(e[3],16));if(e=/#([a-fA-F0-9])([a-fA-F0-9])([a-fA-F0-9])/.exec(t))return i(parseInt(e[1]+e[1],16),parseInt(e[2]+e[2],16),parseInt(e[3]+e[3],16));t=t.trim().toLowerCase();return"transparent"==t?i(255,255,255,0):i((e=o[t]||[0,0,0])[0],e[1],e[2])};var o={aqua:[0,255,255],azure:[240,255,255],beige:[245,245,220],black:[0,0,0],blue:[0,0,255],brown:[165,42,42],cyan:[0,255,255],darkblue:[0,0,139],darkcyan:[0,139,139],darkgrey:[169,169,169],darkgreen:[0,100,0],darkkhaki:[189,183,107],darkmagenta:[139,0,139],darkolivegreen:[85,107,47],darkorange:[255,140,0],darkorchid:[153,50,204],darkred:[139,0,0],darksalmon:[233,150,122],darkviolet:[148,0,211],fuchsia:[255,0,255],gold:[255,215,0],green:[0,128,0],indigo:[75,0,130],khaki:[240,230,140],lightblue:[173,216,230],lightcyan:[224,255,255],lightgreen:[144,238,144],lightgrey:[211,211,211],lightpink:[255,182,193],lightyellow:[255,255,224],lime:[0,255,0],magenta:[255,0,255],maroon:[128,0,0],navy:[0,0,128],olive:[128,128,0],orange:[255,165,0],pink:[255,192,203],purple:[128,0,128],violet:[128,0,128],red:[255,0,0],silver:[192,192,192],white:[255,255,255],yellow:[255,255,0]}}(jQuery),function(_){var d=Object.prototype.hasOwnProperty;function V(t,e){var i=e.children("."+t)[0];if(null==i&&((i=document.createElement("canvas")).className=t,_(i).css({direction:"ltr",position:"absolute",left:0,top:0}).appendTo(e),!i.getContext)){if(!window.G_vmlCanvasManager)throw new Error("Canvas is not available. If you're using IE with a fall-back such as Excanvas, then there's either a mistake in your conditional include, or the page has no DOCTYPE and is rendering in Quirks Mode.");i=window.G_vmlCanvasManager.initElement(i)}this.element=i;t=this.context=i.getContext("2d"),i=window.devicePixelRatio||1,t=t.webkitBackingStorePixelRatio||t.mozBackingStorePixelRatio||t.msBackingStorePixelRatio||t.oBackingStorePixelRatio||t.backingStorePixelRatio||1;this.pixelRatio=i/t,this.resize(e.width(),e.height()),this.textContainer=null,this.text={},this._textCache={}}function o(d,t,e,o){var T=[],M={colors:["#edc240","#afd8f8","#cb4b4b","#4da74d","#9440ed"],legend:{show:!0,noColumns:1,labelFormatter:null,labelBoxBorderColor:"#ccc",container:null,position:"ne",margin:5,backgroundColor:null,backgroundOpacity:.85,sorted:null},xaxis:{show:null,position:"bottom",mode:null,font:null,color:null,tickColor:null,transform:null,inverseTransform:null,min:null,max:null,autoscaleMargin:null,ticks:null,tickFormatter:null,labelWidth:null,labelHeight:null,reserveSpace:null,tickLength:null,alignTicksWithAxis:null,tickDecimals:null,tickSize:null,minTickSize:null},yaxis:{autoscaleMargin:.02,position:"left"},xaxes:[],yaxes:[],series:{points:{show:!1,radius:3,lineWidth:2,fill:!0,fillColor:"#ffffff",symbol:"circle"},lines:{lineWidth:2,fill:!1,fillColor:null,steps:!1},bars:{show:!1,lineWidth:2,barWidth:1,fill:!0,fillColor:null,align:"left",horizontal:!1,zero:!0},shadowSize:3,highlightColor:null},grid:{show:!0,aboveData:!1,color:"#545454",backgroundColor:null,borderColor:null,tickColor:null,margin:0,labelMargin:5,axisMargin:8,borderWidth:2,minBorderMargin:null,markings:null,markingsColor:"#f4f4f4",markingsLineWidth:2,clickable:!1,hoverable:!1,autoHighlight:!0,mouseActiveRadius:10},interaction:{redrawOverlayInterval:1e3/60},hooks:{}},u=null,s=null,h=null,k=null,c=null,p=[],m=[],y={left:0,right:0,top:0,bottom:0},w=0,C=0,S={processOptions:[],processRawData:[],processDatapoints:[],processOffset:[],drawBackground:[],drawSeries:[],draw:[],bindEvents:[],drawOverlay:[],shutdown:[]},W=this;function z(t,e){e=[W].concat(e);for(var i=0;i<t.length;++i)t[i].apply(this,e)}function i(t){T=function(t){for(var e=[],i=0;i<t.length;++i){var o=_.extend(!0,{},M.series);null!=t[i].data?(o.data=t[i].data,delete t[i].data,_.extend(!0,o,t[i]),t[i].data=o.data):o.data=t[i],e.push(o)}return e}(t),function(){var t,e=T.length,i=-1;for(t=0;t<T.length;++t){var o=T[t].color;null!=o&&(e--,"number"==typeof o&&i<o&&(i=o))}e<=i&&(e=i+1);var n,a=[],r=M.colors,l=r.length,s=0;for(t=0;t<e;t++)n=_.color.parse(r[t%l]||"#666"),t%l==0&&t&&(s=0<=s?s<.5?-s-.2:0:-s),a[t]=n.scale("rgb",1+s);var c,h=0;for(t=0;t<T.length;++t){if(null==(c=T[t]).color?(c.color=a[h].toString(),++h):"number"==typeof c.color&&(c.color=a[c.color].toString()),null==c.lines.show){var f,u=!0;for(f in c)if(c[f]&&c[f].show){u=!1;break}u&&(c.lines.show=!0)}null==c.lines.zero&&(c.lines.zero=!!c.lines.fill),c.xaxis=g(p,x(c,"x")),c.yaxis=g(m,x(c,"y"))}}(),function(){var t,e,i,o,n,a,r,l,s,c,h,f,u,d,p=Number.POSITIVE_INFINITY,m=Number.NEGATIVE_INFINITY,x=Number.MAX_VALUE;function g(t,e,i){e<t.datamin&&e!=-x&&(t.datamin=e),i>t.datamax&&i!=x&&(t.datamax=i)}for(_.each(I(),function(t,e){e.datamin=p,e.datamax=m,e.used=!1}),t=0;t<T.length;++t)(n=T[t]).datapoints={points:[]},z(S.processRawData,[n,n.data,n.datapoints]);for(t=0;t<T.length;++t)if(n=T[t],h=n.data,(f=n.datapoints.format)||((f=[]).push({x:!0,number:!0,required:!0}),f.push({y:!0,number:!0,required:!0}),(n.bars.show||n.lines.show&&n.lines.fill)&&(d=!!(n.bars.show&&n.bars.zero||n.lines.show&&n.lines.zero),f.push({y:!0,number:!0,required:!1,defaultValue:0,autoscale:d}),n.bars.horizontal&&(delete f[f.length-1].y,f[f.length-1].x=!0)),n.datapoints.format=f),null==n.datapoints.pointsize){n.datapoints.pointsize=f.length,r=n.datapoints.pointsize,a=n.datapoints.points;var b=n.lines.show&&n.lines.steps;for(n.xaxis.used=n.yaxis.used=!0,e=i=0;e<h.length;++e,i+=r){var v=null==(c=h[e]);if(!v)for(o=0;o<r;++o)l=c[o],(s=f[o])&&(s.number&&null!=l&&(l=+l,isNaN(l)?l=null:l==1/0?l=x:l==-1/0&&(l=-x)),null==l&&(s.required&&(v=!0),null!=s.defaultValue&&(l=s.defaultValue))),a[i+o]=l;if(v)for(o=0;o<r;++o)null!=(l=a[i+o])&&!1!==(s=f[o]).autoscale&&(s.x&&g(n.xaxis,l,l),s.y&&g(n.yaxis,l,l)),a[i+o]=null;else if(b&&0<i&&null!=a[i-r]&&a[i-r]!=a[i]&&a[i-r+1]!=a[i+1]){for(o=0;o<r;++o)a[i+r+o]=a[i+o];a[i+1]=a[i-r+1],i+=r}}}for(t=0;t<T.length;++t)n=T[t],z(S.processDatapoints,[n,n.datapoints]);for(t=0;t<T.length;++t){n=T[t],a=n.datapoints.points,r=n.datapoints.pointsize,f=n.datapoints.format;var k=p,y=p,w=m,M=m;for(e=0;e<a.length;e+=r)if(null!=a[e])for(o=0;o<r;++o)l=a[e+o],(s=f[o])&&!1!==s.autoscale&&l!=x&&l!=-x&&(s.x&&(l<k&&(k=l),w<l&&(w=l)),s.y&&(l<y&&(y=l),M<l&&(M=l)));if(n.bars.show){switch(n.bars.align){case"left":u=0;break;case"right":u=-n.bars.barWidth;break;default:u=-n.bars.barWidth/2}n.bars.horizontal?(y+=u,M+=u+n.bars.barWidth):(k+=u,w+=u+n.bars.barWidth)}g(n.xaxis,k,w),g(n.yaxis,y,M)}_.each(I(),function(t,e){e.datamin==p&&(e.datamin=null),e.datamax==m&&(e.datamax=null)})}()}function x(t,e){e=t[e+"axis"];return"object"==typeof e&&(e=e.n),"number"!=typeof e&&(e=1),e}function I(){return _.grep(p.concat(m),function(t){return t})}function f(t){for(var e,i={},o=0;o<p.length;++o)(e=p[o])&&e.used&&(i["x"+e.n]=e.c2p(t.left));for(o=0;o<m.length;++o)(e=m[o])&&e.used&&(i["y"+e.n]=e.c2p(t.top));return void 0!==i.x1&&(i.x=i.x1),void 0!==i.y1&&(i.y=i.y1),i}function g(t,e){return t[e-1]||(t[e-1]={n:e,direction:t==p?"x":"y",options:_.extend(!0,{},t==p?M.xaxis:M.yaxis)}),t[e-1]}function n(){N&&clearTimeout(N),h.off("mousemove",F),h.off("mouseleave",D),h.off("click",L),z(S.shutdown,[h])}function a(){var t,e,i=I(),o=M.grid.show;for(e in y){var n=M.grid.margin||0;y[e]="number"==typeof n?n:n[e]||0}for(e in z(S.processOffset,[y]),y)"object"==typeof M.grid.borderWidth?y[e]+=o?M.grid.borderWidth[e]:0:y[e]+=o?M.grid.borderWidth:0;if(_.each(i,function(t,e){var i,o,n,a,r=e.options;e.show=null==r.show?e.used:r.show,e.reserveSpace=null==r.reserveSpace?e.show:r.reserveSpace,n=(i=e).options,a=+(null!=n.min?n.min:i.datamin),r=+(null!=n.max?n.max:i.datamax),0==(e=r-a)?(o=0==r?1:.01,null==n.min&&(a-=o),null!=n.max&&null==n.min||(r+=o)):null!=(o=n.autoscaleMargin)&&(null==n.min&&(a-=e*o)<0&&null!=i.datamin&&0<=i.datamin&&(a=0),null==n.max&&0<(r+=e*o)&&null!=i.datamax&&i.datamax<=0&&(r=0)),i.min=a,i.max=r}),o){var a=_.grep(i,function(t){return t.show||t.reserveSpace});for(_.each(a,function(t,e){var i,o;!function(t){var i=t.options;s="number"==typeof i.ticks&&0<i.ticks?i.ticks:.3*Math.sqrt("x"==t.direction?u.width:u.height);var e=(t.max-t.min)/s,o=-Math.floor(Math.log(e)/Math.LN10),n=i.tickDecimals;null!=n&&n<o&&(o=n);var a,r,l=Math.pow(10,-o),s=e/l;if(s<1.5?r=1:s<3?(r=2,2.25<s&&(null==n||o+1<=n)&&(r=2.5,++o)):r=s<7.5?5:10,r*=l,null!=i.minTickSize&&r<i.minTickSize&&(r=i.minTickSize),t.delta=e,t.tickDecimals=Math.max(0,null!=n?n:o),t.tickSize=i.tickSize||r,"time"==i.mode&&!t.tickGenerator)throw new Error("Time mode requires the flot.time plugin.");t.tickGenerator||(t.tickGenerator=function(t){for(var e,i,o,n=[],a=(i=t.min,o=t.tickSize,o*Math.floor(i/o)),r=0,l=Number.NaN;e=l,l=a+r*t.tickSize,n.push(l),++r,l<t.max&&l!=e;);return n},t.tickFormatter=function(t,e){var i=e.tickDecimals?Math.pow(10,e.tickDecimals):1,o=""+Math.round(t*i)/i;if(null!=e.tickDecimals){t=o.indexOf("."),t=-1==t?0:o.length-t-1;if(t<e.tickDecimals)return(t?o:o+".")+(""+i).substr(1,e.tickDecimals-t)}return o}),"function"==typeof i.tickFormatter&&(t.tickFormatter=function(t,e){return""+i.tickFormatter(t,e)}),null==i.alignTicksWithAxis||(a=("x"==t.direction?p:m)[i.alignTicksWithAxis-1])&&a.used&&a!=t&&(0<(o=t.tickGenerator(t)).length&&(null==i.min&&(t.min=Math.min(t.min,o[0])),null==i.max&&1<o.length&&(t.max=Math.max(t.max,o[o.length-1]))),t.tickGenerator=function(t){for(var e,i=[],o=0;o<a.ticks.length;++o)e=(a.ticks[o].v-a.min)/(a.max-a.min),e=t.min+e*(t.max-t.min),i.push(e);return i},t.mode||null!=i.tickDecimals||(r=Math.max(0,1-Math.floor(Math.log(t.delta)/Math.LN10)),1<(o=t.tickGenerator(t)).length&&/\..*0$/.test((o[1]-o[0]).toFixed(r))||(t.tickDecimals=r)))}(e),function(t){var e,i,o=t.options.ticks,n=[];for(null==o||"number"==typeof o&&0<o?n=t.tickGenerator(t):o&&(n="function"==typeof o?o(t):o),t.ticks=[],e=0;e<n.length;++e){var a=null,r=n[e];"object"==typeof r?(i=+r[0],1<r.length&&(a=r[1])):i=+r,null==a&&(a=t.tickFormatter(i,t)),isNaN(i)||t.ticks.push({v:i,label:a})}}(e),o=(i=e).ticks,i.options.autoscaleMargin&&0<o.length&&(null==i.options.min&&(i.min=Math.min(i.min,o[0].v)),null==i.options.max&&1<o.length&&(i.max=Math.max(i.max,o[o.length-1].v))),function(t){for(var e=t.options,i=t.ticks||[],o=e.labelWidth||0,n=e.labelHeight||0,a=o||("x"==t.direction?Math.floor(u.width/(i.length||1)):null),r=t.direction+"Axis "+t.direction+t.n+"Axis",l="flot-"+t.direction+"-axis flot-"+t.direction+t.n+"-axis "+r,s=e.font||"flot-tick-label tickLabel",c=0;c<i.length;++c){var h=i[c];h.label&&(h=u.getTextInfo(l,h.label,s,null,a),o=Math.max(o,h.width),n=Math.max(n,h.height))}t.labelWidth=e.labelWidth||o,t.labelHeight=e.labelHeight||n}(e)}),t=a.length-1;0<=t;--t)!function(i){var t=i.labelWidth,e=i.labelHeight,o=i.options.position,n="x"===i.direction,a=i.options.tickLength,r=M.grid.axisMargin,l=M.grid.labelMargin,s=!0,c=!0,h=!0,f=!1;_.each(n?p:m,function(t,e){e&&(e.show||e.reserveSpace)&&(e===i?f=!0:e.options.position===o&&(f?c=!1:s=!1),f||(h=!1))}),c&&(r=0),null==a&&(a=h?"full":5),isNaN(+a)||(l+=+a),n?(e+=l,"bottom"==o?(y.bottom+=e+r,i.box={top:u.height-y.bottom,height:e}):(i.box={top:y.top+r,height:e},y.top+=e+r)):(t+=l,"left"==o?(i.box={left:y.left+r,width:t},y.left+=t+r):(y.right+=t+r,i.box={left:u.width-y.right,width:t})),i.position=o,i.tickLength=a,i.box.padding=l,i.innermost=s}(a[t]);!function(){var t,e=M.grid.minBorderMargin;if(null==e)for(t=e=0;t<T.length;++t)e=Math.max(e,2*(T[t].points.radius+T[t].points.lineWidth/2));var i={left:e,right:e,top:e,bottom:e};_.each(I(),function(t,e){e.reserveSpace&&e.ticks&&e.ticks.length&&("x"===e.direction?(i.left=Math.max(i.left,e.labelWidth/2),i.right=Math.max(i.right,e.labelWidth/2)):(i.bottom=Math.max(i.bottom,e.labelHeight/2),i.top=Math.max(i.top,e.labelHeight/2)))}),y.left=Math.ceil(Math.max(i.left,y.left)),y.right=Math.ceil(Math.max(i.right,y.right)),y.top=Math.ceil(Math.max(i.top,y.top)),y.bottom=Math.ceil(Math.max(i.bottom,y.bottom))}(),_.each(a,function(t,e){"x"==(e=e).direction?(e.box.left=y.left-e.labelWidth/2,e.box.width=u.width-y.left-y.right+e.labelWidth):(e.box.top=y.top-e.labelHeight/2,e.box.height=u.height-y.bottom-y.top+e.labelHeight)})}w=u.width-y.left-y.right,C=u.height-y.bottom-y.top,_.each(i,function(t,e){function i(t){return t}var o,n,a,r;n=(e=e).options.transform||i,a=e.options.inverseTransform,r="x"==e.direction?(o=e.scale=w/Math.abs(n(e.max)-n(e.min)),Math.min(n(e.max),n(e.min))):(o=-(o=e.scale=C/Math.abs(n(e.max)-n(e.min))),Math.max(n(e.max),n(e.min))),e.p2c=n==i?function(t){return(t-r)*o}:function(t){return(n(t)-r)*o},e.c2p=a?function(t){return a(r+t/o)}:function(t){return r+t/o}}),o&&_.each(I(),function(t,e){var i,o,n,a,r,l=e.box,s=e.direction+"Axis "+e.direction+e.n+"Axis",c="flot-"+e.direction+"-axis flot-"+e.direction+e.n+"-axis "+s,h=e.options.font||"flot-tick-label tickLabel";if(u.removeText(c),e.show&&0!=e.ticks.length)for(var f=0;f<e.ticks.length;++f)!(i=e.ticks[f]).label||i.v<e.min||i.v>e.max||("x"==e.direction?(a="center",o=y.left+e.p2c(i.v),"bottom"==e.position?n=l.top+l.padding:(n=l.top+l.height-l.padding,r="bottom")):(r="middle",n=y.top+e.p2c(i.v),"left"==e.position?(o=l.left+l.width-l.padding,a="right"):o=l.left+l.padding),u.addText(c,o,n,i.label,h,null,null,a,r))}),function(){if(null!=M.legend.container?_(M.legend.container).html(""):d.find(".legend").remove(),M.legend.show){for(var t,e,i,o,n,a,r,l=[],s=[],c=!1,h=M.legend.labelFormatter,f=0;f<T.length;++f)(t=T[f]).label&&(e=h?h(t.label,t):t.label)&&s.push({label:e,color:t.color});for(M.legend.sorted&&("function"==typeof M.legend.sorted?s.sort(M.legend.sorted):"reverse"==M.legend.sorted?s.reverse():(i="descending"!=M.legend.sorted,s.sort(function(t,e){return t.label==e.label?0:t.label<e.label!=i?1:-1}))),f=0;f<s.length;++f){var u=s[f];f%M.legend.noColumns==0&&(c&&l.push("</tr>"),l.push("<tr>"),c=!0),l.push('<td class="legendColorBox"><div style="border:1px solid '+M.legend.labelBoxBorderColor+';padding:1px"><div style="width:4px;height:0;border:5px solid '+u.color+';overflow:hidden"></div></div></td><td class="legendLabel">'+u.label+"</td>")}c&&l.push("</tr>"),0!=l.length&&(o='<table style="font-size:smaller;color:'+M.grid.color+'">'+l.join("")+"</table>",null!=M.legend.container?_(M.legend.container).html(o):(n="",a=M.legend.position,null==(r=M.legend.margin)[0]&&(r=[r,r]),"n"==a.charAt(0)?n+="top:"+(r[1]+y.top)+"px;":"s"==a.charAt(0)&&(n+="bottom:"+(r[1]+y.bottom)+"px;"),"e"==a.charAt(1)?n+="right:"+(r[0]+y.right)+"px;":"w"==a.charAt(1)&&(n+="left:"+(r[0]+y.left)+"px;"),a=_('<div class="legend">'+o.replace('style="','style="position:absolute;'+n+";")+"</div>").appendTo(d),0!=M.legend.backgroundOpacity&&(null==(r=M.legend.backgroundColor)&&((r=(r=M.grid.backgroundColor)&&"string"==typeof r?_.color.parse(r):_.color.extract(a,"background-color")).a=1,r=r.toString()),o=a.children(),_('<div style="position:absolute;width:'+o.width()+"px;height:"+o.height()+"px;"+n+"background-color:"+r+';"> </div>').prependTo(a).css("opacity",M.legend.backgroundOpacity))))}}()}function r(){u.clear(),z(S.drawBackground,[k]);var t=M.grid;t.show&&t.backgroundColor&&(k.save(),k.translate(y.left,y.top),k.fillStyle=G(M.grid.backgroundColor,C,0,"rgba(255, 255, 255, 0)"),k.fillRect(0,0,w,C),k.restore()),t.show&&!t.aboveData&&l();for(var e,i=0;i<T.length;++i)z(S.drawSeries,[k,T[i]]),(e=T[i]).lines.show&&function(t){function e(t,e,i,o,n){var a=t.points,r=t.pointsize,l=null,s=null;k.beginPath();for(var c=r;c<a.length;c+=r){var h=a[c-r],f=a[c-r+1],u=a[c],d=a[c+1];if(null!=h&&null!=u){if(f<=d&&f<n.min){if(d<n.min)continue;h=(n.min-f)/(d-f)*(u-h)+h,f=n.min}else if(d<=f&&d<n.min){if(f<n.min)continue;u=(n.min-f)/(d-f)*(u-h)+h,d=n.min}if(d<=f&&f>n.max){if(d>n.max)continue;h=(n.max-f)/(d-f)*(u-h)+h,f=n.max}else if(f<=d&&d>n.max){if(f>n.max)continue;u=(n.max-f)/(d-f)*(u-h)+h,d=n.max}if(h<=u&&h<o.min){if(u<o.min)continue;f=(o.min-h)/(u-h)*(d-f)+f,h=o.min}else if(u<=h&&u<o.min){if(h<o.min)continue;d=(o.min-h)/(u-h)*(d-f)+f,u=o.min}if(u<=h&&h>o.max){if(u>o.max)continue;f=(o.max-h)/(u-h)*(d-f)+f,h=o.max}else if(h<=u&&u>o.max){if(h>o.max)continue;d=(o.max-h)/(u-h)*(d-f)+f,u=o.max}h==l&&f==s||k.moveTo(o.p2c(h)+e,n.p2c(f)+i),l=u,s=d,k.lineTo(o.p2c(u)+e,n.p2c(d)+i)}}k.stroke()}k.save(),k.translate(y.left,y.top),k.lineJoin="round";var i,o=t.lines.lineWidth,n=t.shadowSize;0<o&&0<n&&(k.lineWidth=n,k.strokeStyle="rgba(0,0,0,0.1)",i=Math.PI/18,e(t.datapoints,Math.sin(i)*(o/2+n/2),Math.cos(i)*(o/2+n/2),t.xaxis,t.yaxis),k.lineWidth=n/2,e(t.datapoints,Math.sin(i)*(o/2+n/4),Math.cos(i)*(o/2+n/4),t.xaxis,t.yaxis)),k.lineWidth=o,k.strokeStyle=t.color,(n=v(t.lines,t.color,0,C))&&(k.fillStyle=n,function(t,e,i){for(var o=t.points,n=t.pointsize,a=Math.min(Math.max(0,i.min),i.max),r=0,l=!1,s=1,c=0,h=0;!(0<n&&r>o.length+n);){var f,u,d=o[(r+=n)-n],p=o[r-n+s],m=o[r],x=o[r+s];if(l){if(0<n&&null!=d&&null==m){h=r,n=-n,s=2;continue}if(n<0&&r==c+n){k.fill(),l=!1,s=1,r=c=h+(n=-n);continue}}if(null!=d&&null!=m){if(d<=m&&d<e.min){if(m<e.min)continue;p=(e.min-d)/(m-d)*(x-p)+p,d=e.min}else if(m<=d&&m<e.min){if(d<e.min)continue;x=(e.min-d)/(m-d)*(x-p)+p,m=e.min}if(m<=d&&d>e.max){if(m>e.max)continue;p=(e.max-d)/(m-d)*(x-p)+p,d=e.max}else if(d<=m&&m>e.max){if(d>e.max)continue;x=(e.max-d)/(m-d)*(x-p)+p,m=e.max}l||(k.beginPath(),k.moveTo(e.p2c(d),i.p2c(a)),l=!0),p>=i.max&&x>=i.max?(k.lineTo(e.p2c(d),i.p2c(i.max)),k.lineTo(e.p2c(m),i.p2c(i.max))):p<=i.min&&x<=i.min?(k.lineTo(e.p2c(d),i.p2c(i.min)),k.lineTo(e.p2c(m),i.p2c(i.min))):(f=d,u=m,p<=x&&p<i.min&&x>=i.min?(d=(i.min-p)/(x-p)*(m-d)+d,p=i.min):x<=p&&x<i.min&&p>=i.min&&(m=(i.min-p)/(x-p)*(m-d)+d,x=i.min),x<=p&&p>i.max&&x<=i.max?(d=(i.max-p)/(x-p)*(m-d)+d,p=i.max):p<=x&&x>i.max&&p<=i.max&&(m=(i.max-p)/(x-p)*(m-d)+d,x=i.max),d!=f&&k.lineTo(e.p2c(f),i.p2c(p)),k.lineTo(e.p2c(d),i.p2c(p)),k.lineTo(e.p2c(m),i.p2c(x)),m!=u&&(k.lineTo(e.p2c(m),i.p2c(x)),k.lineTo(e.p2c(u),i.p2c(x))))}}}(t.datapoints,t.xaxis,t.yaxis)),0<o&&e(t.datapoints,0,0,t.xaxis,t.yaxis),k.restore()}(e),e.bars.show&&function(c){var t;switch(k.save(),k.translate(y.left,y.top),k.lineWidth=c.bars.lineWidth,k.strokeStyle=c.color,c.bars.align){case"left":t=0;break;case"right":t=-c.bars.barWidth;break;default:t=-c.bars.barWidth/2}var e=c.bars.fill?function(t,e){return v(c.bars,c.color,t,e)}:null;(function(t,e,i,o,n,a){for(var r=t.points,l=t.pointsize,s=0;s<r.length;s+=l)null!=r[s]&&b(r[s],r[s+1],r[s+2],e,i,o,n,a,k,c.bars.horizontal,c.bars.lineWidth)})(c.datapoints,t,t+c.bars.barWidth,e,c.xaxis,c.yaxis),k.restore()}(e),e.points.show&&function(t){function e(t,e,i,o,n,a,r,l){for(var s=t.points,c=t.pointsize,h=0;h<s.length;h+=c){var f=s[h],u=s[h+1];null==f||f<a.min||f>a.max||u<r.min||u>r.max||(k.beginPath(),f=a.p2c(f),u=r.p2c(u)+o,"circle"==l?k.arc(f,u,e,0,n?Math.PI:2*Math.PI,!1):l(k,f,u,e,n),k.closePath(),i&&(k.fillStyle=i,k.fill()),k.stroke())}}k.save(),k.translate(y.left,y.top);var i=t.points.lineWidth,o=t.shadowSize,n=t.points.radius,a=t.points.symbol;0==i&&(i=1e-4),0<i&&0<o&&(o=o/2,k.lineWidth=o,k.strokeStyle="rgba(0,0,0,0.1)",e(t.datapoints,n,null,o+o/2,!0,t.xaxis,t.yaxis,a),k.strokeStyle="rgba(0,0,0,0.2)",e(t.datapoints,n,null,o/2,!0,t.xaxis,t.yaxis,a)),k.lineWidth=i,k.strokeStyle=t.color,e(t.datapoints,n,v(t.points,t.color),0,!1,t.xaxis,t.yaxis,a),k.restore()}(e);z(S.draw,[k]),t.show&&t.aboveData&&l(),u.render(),R()}function A(t,e){for(var i,o,n,a,r,l=I(),s=0;s<l.length;++s)if(i=l[s],i.direction==e&&(t[o=e+i.n+"axis"]||1!=i.n||(o=e+"axis"),t[o])){a=t[o].from,r=t[o].to;break}return t[o]||(i=("x"==e?p:m)[0],a=t[e+"1"],r=t[e+"2"]),null!=a&&null!=r&&r<a&&(n=a,a=r,r=n),{from:a,to:r,axis:i}}function l(){var t,e,i,o;k.save(),k.translate(y.left,y.top);var n=M.grid.markings;if(n)for("function"==typeof n&&((e=W.getAxes()).xmin=e.xaxis.min,e.xmax=e.xaxis.max,e.ymin=e.yaxis.min,e.ymax=e.yaxis.max,n=n(e)),t=0;t<n.length;++t){var a,r,l,s=n[t],c=A(s,"x"),h=A(s,"y");null==c.from&&(c.from=c.axis.min),null==c.to&&(c.to=c.axis.max),null==h.from&&(h.from=h.axis.min),null==h.to&&(h.to=h.axis.max),c.to<c.axis.min||c.from>c.axis.max||h.to<h.axis.min||h.from>h.axis.max||(c.from=Math.max(c.from,c.axis.min),c.to=Math.min(c.to,c.axis.max),h.from=Math.max(h.from,h.axis.min),h.to=Math.min(h.to,h.axis.max),a=c.from===c.to,l=h.from===h.to,a&&l||(c.from=Math.floor(c.axis.p2c(c.from)),c.to=Math.floor(c.axis.p2c(c.to)),h.from=Math.floor(h.axis.p2c(h.from)),h.to=Math.floor(h.axis.p2c(h.to)),a||l?(l=(r=s.lineWidth||M.grid.markingsLineWidth)%2?.5:0,k.beginPath(),k.strokeStyle=s.color||M.grid.markingsColor,k.lineWidth=r,a?(k.moveTo(c.to+l,h.from),k.lineTo(c.to+l,h.to)):(k.moveTo(c.from,h.to+l),k.lineTo(c.to,h.to+l)),k.stroke()):(k.fillStyle=s.color||M.grid.markingsColor,k.fillRect(c.from,h.to,c.to-c.from,h.from-h.to))))}e=I(),i=M.grid.borderWidth;for(var f=0;f<e.length;++f){var u,d,p,m=e[f],x=m.box,g=m.tickLength;if(m.show&&0!=m.ticks.length){for(k.lineWidth=1,"x"==m.direction?(u=0,d="full"==g?"top"==m.position?0:C:x.top-y.top+("top"==m.position?x.height:0)):(d=0,u="full"==g?"left"==m.position?0:w:x.left-y.left+("left"==m.position?x.width:0)),m.innermost||(k.strokeStyle=m.options.color,k.beginPath(),v=p=0,"x"==m.direction?v=w+1:p=C+1,1==k.lineWidth&&("x"==m.direction?d=Math.floor(d)+.5:u=Math.floor(u)+.5),k.moveTo(u,d),k.lineTo(u+v,d+p),k.stroke()),k.strokeStyle=m.options.tickColor,k.beginPath(),t=0;t<m.ticks.length;++t){var b=m.ticks[t].v,v=p=0;isNaN(b)||b<m.min||b>m.max||"full"==g&&("object"==typeof i&&0<i[m.position]||0<i)&&(b==m.min||b==m.max)||("x"==m.direction?(u=m.p2c(b),p="full"==g?-C:g,"top"==m.position&&(p=-p)):(d=m.p2c(b),v="full"==g?-w:g,"left"==m.position&&(v=-v)),1==k.lineWidth&&("x"==m.direction?u=Math.floor(u)+.5:d=Math.floor(d)+.5),k.moveTo(u,d),k.lineTo(u+v,d+p))}k.stroke()}}i&&(o=M.grid.borderColor,"object"==typeof i||"object"==typeof o?("object"!=typeof i&&(i={top:i,right:i,bottom:i,left:i}),"object"!=typeof o&&(o={top:o,right:o,bottom:o,left:o}),0<i.top&&(k.strokeStyle=o.top,k.lineWidth=i.top,k.beginPath(),k.moveTo(0-i.left,0-i.top/2),k.lineTo(w,0-i.top/2),k.stroke()),0<i.right&&(k.strokeStyle=o.right,k.lineWidth=i.right,k.beginPath(),k.moveTo(w+i.right/2,0-i.top),k.lineTo(w+i.right/2,C),k.stroke()),0<i.bottom&&(k.strokeStyle=o.bottom,k.lineWidth=i.bottom,k.beginPath(),k.moveTo(w+i.right,C+i.bottom/2),k.lineTo(0,C+i.bottom/2),k.stroke()),0<i.left&&(k.strokeStyle=o.left,k.lineWidth=i.left,k.beginPath(),k.moveTo(0-i.left/2,C+i.bottom),k.lineTo(0-i.left/2,0),k.stroke())):(k.lineWidth=i,k.strokeStyle=M.grid.borderColor,k.strokeRect(-i/2,-i/2,w+i,C+i))),k.restore()}function b(t,e,i,o,n,a,r,l,s,c,h){var f,u,d,p,m,x,g,b,v;c?(m=!(b=x=g=!0),p=e+o,d=e+n,(u=t)<(f=i)&&(v=u,u=f,f=v,x=!(m=!0))):(b=!(m=x=g=!0),f=t+o,u=t+n,(p=e)<(d=i)&&(v=p,p=d,d=v,g=!(b=!0))),u<r.min||f>r.max||p<l.min||d>l.max||(f<r.min&&(f=r.min,m=!1),u>r.max&&(u=r.max,x=!1),d<l.min&&(d=l.min,b=!1),p>l.max&&(p=l.max,g=!1),f=r.p2c(f),d=l.p2c(d),u=r.p2c(u),p=l.p2c(p),a&&(s.fillStyle=a(d,p),s.fillRect(f,p,u-f,d-p)),0<h&&(m||x||g||b)&&(s.beginPath(),s.moveTo(f,d),m?s.lineTo(f,p):s.moveTo(f,p),g?s.lineTo(u,p):s.moveTo(u,p),x?s.lineTo(u,d):s.moveTo(u,d),b?s.lineTo(f,d):s.moveTo(f,d),s.stroke()))}function v(t,e,i,o){var n=t.fill;if(!n)return null;if(t.fillColor)return G(t.fillColor,i,o,e);e=_.color.parse(e);return e.a="number"==typeof n?n:.4,e.normalize(),e.toString()}W.setData=i,W.setupGrid=a,W.draw=r,W.getPlaceholder=function(){return d},W.getCanvas=function(){return u.element},W.getPlotOffset=function(){return y},W.width=function(){return w},W.height=function(){return C},W.offset=function(){var t=h.offset();return t.left+=y.left,t.top+=y.top,t},W.getData=function(){return T},W.getAxes=function(){var i={};return _.each(p.concat(m),function(t,e){e&&(i[e.direction+(1!=e.n?e.n:"")+"axis"]=e)}),i},W.getXAxes=function(){return p},W.getYAxes=function(){return m},W.c2p=f,W.p2c=function(t){var e,i,o,n={};for(e=0;e<p.length;++e)if(i=p[e],i&&i.used&&(o="x"+i.n,null==t[o]&&1==i.n&&(o="x"),null!=t[o])){n.left=i.p2c(t[o]);break}for(e=0;e<m.length;++e)if(i=m[e],i&&i.used&&(o="y"+i.n,null==t[o]&&1==i.n&&(o="y"),null!=t[o])){n.top=i.p2c(t[o]);break}return n},W.getOptions=function(){return M},W.highlight=j,W.unhighlight=E,W.triggerRedrawOverlay=R,W.pointOffset=function(t){return{left:parseInt(p[x(t,"x")-1].p2c(+t.x)+y.left,10),top:parseInt(m[x(t,"y")-1].p2c(+t.y)+y.top,10)}},W.shutdown=n,W.destroy=function(){n(),d.removeData("plot").empty(),T=[],p=[],m=[],P=[],W=S=c=k=h=s=u=M=null},W.resize=function(){var t=d.width(),e=d.height();u.resize(t,e),s.resize(t,e)},W.hooks=S,function(){for(var t={Canvas:V},e=0;e<o.length;++e){var i=o[e];i.init(W,t),i.options&&_.extend(!0,M,i.options)}}(),function(t){_.extend(!0,M,t),t&&t.colors&&(M.colors=t.colors),null==M.xaxis.color&&(M.xaxis.color=_.color.parse(M.grid.color).scale("a",.22).toString()),null==M.yaxis.color&&(M.yaxis.color=_.color.parse(M.grid.color).scale("a",.22).toString()),null==M.xaxis.tickColor&&(M.xaxis.tickColor=M.grid.tickColor||M.xaxis.color),null==M.yaxis.tickColor&&(M.yaxis.tickColor=M.grid.tickColor||M.yaxis.color),null==M.grid.borderColor&&(M.grid.borderColor=M.grid.color),null==M.grid.tickColor&&(M.grid.tickColor=_.color.parse(M.grid.color).scale("a",.22).toString());var e,i,o,n,t=(t=d.css("font-size"))?+t.replace("px",""):13,a={style:d.css("font-style"),size:Math.round(.8*t),variant:d.css("font-variant"),weight:d.css("font-weight"),family:d.css("font-family")};for(o=M.xaxes.length||1,e=0;e<o;++e)(i=M.xaxes[e])&&!i.tickColor&&(i.tickColor=i.color),i=_.extend(!0,{},M.xaxis,i),(M.xaxes[e]=i).font&&(i.font=_.extend({},a,i.font),i.font.color||(i.font.color=i.color),i.font.lineHeight||(i.font.lineHeight=Math.round(1.15*i.font.size)));for(o=M.yaxes.length||1,e=0;e<o;++e)(i=M.yaxes[e])&&!i.tickColor&&(i.tickColor=i.color),i=_.extend(!0,{},M.yaxis,i),(M.yaxes[e]=i).font&&(i.font=_.extend({},a,i.font),i.font.color||(i.font.color=i.color),i.font.lineHeight||(i.font.lineHeight=Math.round(1.15*i.font.size)));for(M.xaxis.noTicks&&null==M.xaxis.ticks&&(M.xaxis.ticks=M.xaxis.noTicks),M.yaxis.noTicks&&null==M.yaxis.ticks&&(M.yaxis.ticks=M.yaxis.noTicks),M.x2axis&&(M.xaxes[1]=_.extend(!0,{},M.xaxis,M.x2axis),M.xaxes[1].position="top",null==M.x2axis.min&&(M.xaxes[1].min=null),null==M.x2axis.max&&(M.xaxes[1].max=null)),M.y2axis&&(M.yaxes[1]=_.extend(!0,{},M.yaxis,M.y2axis),M.yaxes[1].position="right",null==M.y2axis.min&&(M.yaxes[1].min=null),null==M.y2axis.max&&(M.yaxes[1].max=null)),M.grid.coloredAreas&&(M.grid.markings=M.grid.coloredAreas),M.grid.coloredAreasColor&&(M.grid.markingsColor=M.grid.coloredAreasColor),M.lines&&_.extend(!0,M.series.lines,M.lines),M.points&&_.extend(!0,M.series.points,M.points),M.bars&&_.extend(!0,M.series.bars,M.bars),null!=M.shadowSize&&(M.series.shadowSize=M.shadowSize),null!=M.highlightColor&&(M.series.highlightColor=M.highlightColor),e=0;e<M.xaxes.length;++e)g(p,e+1).options=M.xaxes[e];for(e=0;e<M.yaxes.length;++e)g(m,e+1).options=M.yaxes[e];for(n in S)M.hooks[n]&&M.hooks[n].length&&(S[n]=S[n].concat(M.hooks[n]));z(S.processOptions,[M])}(e),function(){d.css("padding",0).children().filter(function(){return!_(this).hasClass("flot-overlay")&&!_(this).hasClass("flot-base")}).remove(),"static"==d.css("position")&&d.css("position","relative"),u=new V("flot-base",d),s=new V("flot-overlay",d),k=u.context,c=s.context,h=_(s.element).off();var t=d.data("plot");t&&(t.shutdown(),s.clear()),d.data("plot",W)}(),i(t),a(),r(),M.grid.hoverable&&(h.on("mousemove",F),h.on("mouseleave",D)),M.grid.clickable&&h.on("click",L),z(S.bindEvents,[h]);var P=[],N=null;function F(t){M.grid.hoverable&&O("plothover",t,function(t){return 0!=t.hoverable})}function D(t){M.grid.hoverable&&O("plothover",t,function(t){return!1})}function L(t){O("plotclick",t,function(t){return 0!=t.clickable})}function O(t,e,i){var o=h.offset(),n=e.pageX-o.left-y.left,a=e.pageY-o.top-y.top,r=f({left:n,top:a});r.pageX=e.pageX,r.pageY=e.pageY;var l=function(t,e,i){for(var o,n=M.grid.mouseActiveRadius,a=n*n+1,r=null,l=T.length-1;0<=l;--l)if(i(T[l])){var s,c,h=T[l],f=h.xaxis,u=h.yaxis,d=h.datapoints.points,p=f.c2p(t),m=u.c2p(e),x=n/f.scale,g=n/u.scale,b=h.datapoints.pointsize;if(f.options.inverseTransform&&(x=Number.MAX_VALUE),u.options.inverseTransform&&(g=Number.MAX_VALUE),h.lines.show||h.points.show)for(o=0;o<d.length;o+=b){var v,k=d[o],y=d[o+1];null!=k&&(x<k-p||k-p<-x||g<y-m||y-m<-g||(v=(v=Math.abs(f.p2c(k)-t))*v+(v=Math.abs(u.p2c(y)-e))*v)<a&&(a=v,r=[l,o/b]))}if(h.bars.show&&!r){switch(h.bars.align){case"left":s=0;break;case"right":s=-h.bars.barWidth;break;default:s=-h.bars.barWidth/2}for(c=s+h.bars.barWidth,o=0;o<d.length;o+=b){var k=d[o],y=d[o+1],w=d[o+2];null!=k&&(T[l].bars.horizontal?p<=Math.max(w,k)&&p>=Math.min(w,k)&&y+s<=m&&m<=y+c:k+s<=p&&p<=k+c&&m>=Math.min(w,y)&&m<=Math.max(w,y))&&(r=[l,o/b])}}}return r?(l=r[0],o=r[1],b=T[l].datapoints.pointsize,{datapoint:T[l].datapoints.points.slice(o*b,(o+1)*b),dataIndex:o,series:T[l],seriesIndex:l}):null}(n,a,i);if(l&&(l.pageX=parseInt(l.series.xaxis.p2c(l.datapoint[0])+o.left+y.left,10),l.pageY=parseInt(l.series.yaxis.p2c(l.datapoint[1])+o.top+y.top,10)),M.grid.autoHighlight){for(var s=0;s<P.length;++s){var c=P[s];c.auto!=t||l&&c.series==l.series&&c.point[0]==l.datapoint[0]&&c.point[1]==l.datapoint[1]||E(c.series,c.point)}l&&j(l.series,l.datapoint,t)}d.trigger(t,[r,l])}function R(){var t=M.interaction.redrawOverlayInterval;-1!=t?N=N||setTimeout(H,t):H()}function H(){var t,e,i,o,n,a,r,l;for(N=null,c.save(),s.clear(),c.translate(y.left,y.top),t=0;t<P.length;++t)(l=P[t]).series.bars.show?function(t,e){var i,o="string"==typeof t.highlightColor?t.highlightColor:_.color.parse(t.color).scale("a",.5).toString(),n=o;switch(t.bars.align){case"left":i=0;break;case"right":i=-t.bars.barWidth;break;default:i=-t.bars.barWidth/2}c.lineWidth=t.bars.lineWidth,c.strokeStyle=o,b(e[0],e[1],e[2]||0,i,i+t.bars.barWidth,function(){return n},t.xaxis,t.yaxis,c,t.bars.horizontal,t.bars.lineWidth)}(l.series,l.point):(e=l.series,i=l.point,l=r=a=n=o=void 0,o=i[0],n=i[1],a=e.xaxis,r=e.yaxis,l="string"==typeof e.highlightColor?e.highlightColor:_.color.parse(e.color).scale("a",.5).toString(),o<a.min||o>a.max||n<r.min||n>r.max||(i=e.points.radius+e.points.lineWidth/2,c.lineWidth=i,c.strokeStyle=l,i*=1.5,o=a.p2c(o),n=r.p2c(n),c.beginPath(),"circle"==e.points.symbol?c.arc(o,n,i,0,2*Math.PI,!1):e.points.symbol(c,o,n,i,!1),c.closePath(),c.stroke()));c.restore(),z(S.drawOverlay,[c])}function j(t,e,i){"number"==typeof t&&(t=T[t]),"number"==typeof e&&(o=t.datapoints.pointsize,e=t.datapoints.points.slice(o*e,o*(e+1)));var o=B(t,e);-1==o?(P.push({series:t,point:e,auto:i}),R()):i||(P[o].auto=!1)}function E(t,e){if(null==t&&null==e)return P=[],void R();var i;"number"==typeof t&&(t=T[t]),"number"==typeof e&&(i=t.datapoints.pointsize,e=t.datapoints.points.slice(i*e,i*(e+1)));e=B(t,e);-1!=e&&(P.splice(e,1),R())}function B(t,e){for(var i=0;i<P.length;++i){var o=P[i];if(o.series==t&&o.point[0]==e[0]&&o.point[1]==e[1])return i}return-1}function G(t,e,i,o){if("string"==typeof t)return t;for(var n=k.createLinearGradient(0,i,0,e),a=0,r=t.colors.length;a<r;++a){var l,s=t.colors[a];"string"!=typeof s&&(l=_.color.parse(o),null!=s.brightness&&(l=l.scale("rgb",s.brightness)),null!=s.opacity&&(l.a*=s.opacity),s=l.toString()),n.addColorStop(a/(r-1),s)}return n}}_.fn.detach||(_.fn.detach=function(){return this.each(function(){this.parentNode&&this.parentNode.removeChild(this)})}),V.prototype.resize=function(t,e){if(t<=0||e<=0)throw new Error("Invalid dimensions for plot, width = "+t+", height = "+e);var i=this.element,o=this.context,n=this.pixelRatio;this.width!=t&&(i.width=t*n,i.style.width=t+"px",this.width=t),this.height!=e&&(i.height=e*n,i.style.height=e+"px",this.height=e),o.restore(),o.save(),o.scale(n,n)},V.prototype.clear=function(){this.context.clearRect(0,0,this.width,this.height)},V.prototype.render=function(){var t,e=this._textCache;for(t in e)if(d.call(e,t)){var i,o=this.getTextLayer(t),n=e[t];for(i in o.hide(),n)if(d.call(n,i)){var a,r=n[i];for(a in r)if(d.call(r,a)){for(var l,s=r[a].positions,c=0;l=s[c];c++)l.active?l.rendered||(o.append(l.element),l.rendered=!0):(s.splice(c--,1),l.rendered&&l.element.detach());0==s.length&&delete r[a]}}o.show()}},V.prototype.getTextLayer=function(t){var e=this.text[t];return null==e&&(null==this.textContainer&&(this.textContainer=_("<div class='flot-text'></div>").css({position:"absolute",top:0,left:0,bottom:0,right:0,"font-size":"smaller",color:"#545454"}).insertAfter(this.element)),e=this.text[t]=_("<div></div>").addClass(t).css({position:"absolute",top:0,left:0,bottom:0,right:0}).appendTo(this.textContainer)),e},V.prototype.getTextInfo=function(t,e,i,o,n){var a,r,l;return e=""+e,a="object"==typeof i?i.style+" "+i.variant+" "+i.weight+" "+i.size+"px/"+i.lineHeight+"px "+i.family:i,null==(l=this._textCache[t])&&(l=this._textCache[t]={}),null==(r=l[a])&&(r=l[a]={}),null==(l=r[e])&&(t=_("<div></div>").html(e).css({position:"absolute","max-width":n,top:-9999}).appendTo(this.getTextLayer(t)),"object"==typeof i?t.css({font:a,color:i.color}):"string"==typeof i&&t.addClass(i),l=r[e]={width:t.outerWidth(!0),height:t.outerHeight(!0),element:t,positions:[]},t.detach()),l},V.prototype.addText=function(t,e,i,o,n,a,r,l,s){var r=this.getTextInfo(t,o,n,a,r),c=r.positions;"center"==l?e-=r.width/2:"right"==l&&(e-=r.width),"middle"==s?i-=r.height/2:"bottom"==s&&(i-=r.height);for(var h,f=0;h=c[f];f++)if(h.x==e&&h.y==i)return void(h.active=!0);h={active:!0,rendered:!1,element:c.length?r.element.clone():r.element,x:e,y:i},c.push(h),h.element.css({top:Math.round(i),left:Math.round(e),"text-align":l})},V.prototype.removeText=function(t,e,i,o,n,a){if(null==o){var r=this._textCache[t];if(null!=r)for(var l in r)if(d.call(r,l)){var s,c=r[l];for(s in c)if(d.call(c,s))for(var h=c[s].positions,f=0;u=h[f];f++)u.active=!1}}else for(var u,h=this.getTextInfo(t,o,n,a).positions,f=0;u=h[f];f++)u.x==e&&u.y==i&&(u.active=!1)},_.plot=function(t,e,i){return new o(_(t),e,i,_.plot.plugins)},_.plot.version="0.8.3",_.plot.plugins=[],_.fn.plot=function(t,e){return this.each(function(){_.plot(this,t,e)})}}(jQuery);