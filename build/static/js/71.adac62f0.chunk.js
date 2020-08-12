/*! For license information please see 71.adac62f0.chunk.js.LICENSE.txt */
(this.webpackJsonpgitconvex=this.webpackJsonpgitconvex||[]).push([[71],{148:function(e,n,t){(function(n){var t=function(e){var n=/\blang(?:uage)?-([\w-]+)\b/i,t=0,r={manual:e.Prism&&e.Prism.manual,disableWorkerMessageHandler:e.Prism&&e.Prism.disableWorkerMessageHandler,util:{encode:function e(n){return n instanceof a?new a(n.type,e(n.content),n.alias):Array.isArray(n)?n.map(e):n.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/\u00a0/g," ")},type:function(e){return Object.prototype.toString.call(e).slice(8,-1)},objId:function(e){return e.__id||Object.defineProperty(e,"__id",{value:++t}),e.__id},clone:function e(n,t){var a,i;switch(t=t||{},r.util.type(n)){case"Object":if(i=r.util.objId(n),t[i])return t[i];for(var l in a={},t[i]=a,n)n.hasOwnProperty(l)&&(a[l]=e(n[l],t));return a;case"Array":return i=r.util.objId(n),t[i]?t[i]:(a=[],t[i]=a,n.forEach((function(n,r){a[r]=e(n,t)})),a);default:return n}},getLanguage:function(e){for(;e&&!n.test(e.className);)e=e.parentElement;return e?(e.className.match(n)||[,"none"])[1].toLowerCase():"none"},currentScript:function(){if("undefined"===typeof document)return null;if("currentScript"in document)return document.currentScript;try{throw new Error}catch(r){var e=(/at [^(\r\n]*\((.*):.+:.+\)$/i.exec(r.stack)||[])[1];if(e){var n=document.getElementsByTagName("script");for(var t in n)if(n[t].src==e)return n[t]}return null}},isActive:function(e,n,t){for(var r="no-"+n;e;){var a=e.classList;if(a.contains(n))return!0;if(a.contains(r))return!1;e=e.parentElement}return!!t}},languages:{extend:function(e,n){var t=r.util.clone(r.languages[e]);for(var a in n)t[a]=n[a];return t},insertBefore:function(e,n,t,a){var i=(a=a||r.languages)[e],l={};for(var o in i)if(i.hasOwnProperty(o)){if(o==n)for(var s in t)t.hasOwnProperty(s)&&(l[s]=t[s]);t.hasOwnProperty(o)||(l[o]=i[o])}var u=a[e];return a[e]=l,r.languages.DFS(r.languages,(function(n,t){t===u&&n!=e&&(this[n]=l)})),l},DFS:function e(n,t,a,i){i=i||{};var l=r.util.objId;for(var o in n)if(n.hasOwnProperty(o)){t.call(n,o,n[o],a||o);var s=n[o],u=r.util.type(s);"Object"!==u||i[l(s)]?"Array"!==u||i[l(s)]||(i[l(s)]=!0,e(s,t,o,i)):(i[l(s)]=!0,e(s,t,null,i))}}},plugins:{},highlightAll:function(e,n){r.highlightAllUnder(document,e,n)},highlightAllUnder:function(e,n,t){var a={callback:t,container:e,selector:'code[class*="language-"], [class*="language-"] code, code[class*="lang-"], [class*="lang-"] code'};r.hooks.run("before-highlightall",a),a.elements=Array.prototype.slice.apply(a.container.querySelectorAll(a.selector)),r.hooks.run("before-all-elements-highlight",a);for(var i,l=0;i=a.elements[l++];)r.highlightElement(i,!0===n,a.callback)},highlightElement:function(t,a,i){var l=r.util.getLanguage(t),o=r.languages[l];t.className=t.className.replace(n,"").replace(/\s+/g," ")+" language-"+l;var s=t.parentElement;s&&"pre"===s.nodeName.toLowerCase()&&(s.className=s.className.replace(n,"").replace(/\s+/g," ")+" language-"+l);var u={element:t,language:l,grammar:o,code:t.textContent};function c(e){u.highlightedCode=e,r.hooks.run("before-insert",u),u.element.innerHTML=u.highlightedCode,r.hooks.run("after-highlight",u),r.hooks.run("complete",u),i&&i.call(u.element)}if(r.hooks.run("before-sanity-check",u),!u.code)return r.hooks.run("complete",u),void(i&&i.call(u.element));if(r.hooks.run("before-highlight",u),u.grammar)if(a&&e.Worker){var g=new Worker(r.filename);g.onmessage=function(e){c(e.data)},g.postMessage(JSON.stringify({language:u.language,code:u.code,immediateClose:!0}))}else c(r.highlight(u.code,u.grammar,u.language));else c(r.util.encode(u.code))},highlight:function(e,n,t){var i={code:e,grammar:n,language:t};return r.hooks.run("before-tokenize",i),i.tokens=r.tokenize(i.code,i.grammar),r.hooks.run("after-tokenize",i),a.stringify(r.util.encode(i.tokens),i.language)},tokenize:function(e,n){var t=n.rest;if(t){for(var s in t)n[s]=t[s];delete n.rest}var u=new i;return l(u,u.head,e),function e(n,t,i,s,u,c){for(var g in i)if(i.hasOwnProperty(g)&&i[g]){var h=i[g];h=Array.isArray(h)?h:[h];for(var f=0;f<h.length;++f){if(c&&c.cause==g+","+f)return;var d=h[f],p=d.inside,v=!!d.lookbehind,m=!!d.greedy,y=0,k=d.alias;if(m&&!d.pattern.global){var b=d.pattern.toString().match(/[imsuy]*$/)[0];d.pattern=RegExp(d.pattern.source,b+"g")}for(var x=d.pattern||d,w=s.next,A=u;w!==t.tail&&!(c&&A>=c.reach);A+=w.value.length,w=w.next){var S=w.value;if(t.length>n.length)return;if(!(S instanceof a)){var E=1;if(m&&w!=t.tail.prev){if(x.lastIndex=A,!(j=x.exec(n)))break;var O=j.index+(v&&j[1]?j[1].length:0),P=j.index+j[0].length,L=A;for(L+=w.value.length;O>=L;)w=w.next,L+=w.value.length;if(L-=w.value.length,A=L,w.value instanceof a)continue;for(var N=w;N!==t.tail&&(L<P||"string"===typeof N.value);N=N.next)E++,L+=N.value.length;E--,S=n.slice(A,L),j.index-=A}else{x.lastIndex=0;var j=x.exec(S)}if(j){v&&(y=j[1]?j[1].length:0);O=j.index+y;var C=j[0].slice(y),M=(P=O+C.length,S.slice(0,O)),W=S.slice(P),I=A+S.length;c&&I>c.reach&&(c.reach=I);var _=w.prev;M&&(_=l(t,_,M),A+=M.length),o(t,_,E);var z=new a(g,p?r.tokenize(C,p):C,k,C);w=l(t,_,z),W&&l(t,w,W),E>1&&e(n,t,i,w.prev,A,{cause:g+","+f,reach:I})}}}}}}(e,u,n,u.head,0),function(e){var n=[],t=e.head.next;for(;t!==e.tail;)n.push(t.value),t=t.next;return n}(u)},hooks:{all:{},add:function(e,n){var t=r.hooks.all;t[e]=t[e]||[],t[e].push(n)},run:function(e,n){var t=r.hooks.all[e];if(t&&t.length)for(var a,i=0;a=t[i++];)a(n)}},Token:a};function a(e,n,t,r){this.type=e,this.content=n,this.alias=t,this.length=0|(r||"").length}function i(){var e={value:null,prev:null,next:null},n={value:null,prev:e,next:null};e.next=n,this.head=e,this.tail=n,this.length=0}function l(e,n,t){var r=n.next,a={value:t,prev:n,next:r};return n.next=a,r.prev=a,e.length++,a}function o(e,n,t){for(var r=n.next,a=0;a<t&&r!==e.tail;a++)r=r.next;n.next=r,r.prev=n,e.length-=a}if(e.Prism=r,a.stringify=function e(n,t){if("string"==typeof n)return n;if(Array.isArray(n)){var a="";return n.forEach((function(n){a+=e(n,t)})),a}var i={type:n.type,content:e(n.content,t),tag:"span",classes:["token",n.type],attributes:{},language:t},l=n.alias;l&&(Array.isArray(l)?Array.prototype.push.apply(i.classes,l):i.classes.push(l)),r.hooks.run("wrap",i);var o="";for(var s in i.attributes)o+=" "+s+'="'+(i.attributes[s]||"").replace(/"/g,"&quot;")+'"';return"<"+i.tag+' class="'+i.classes.join(" ")+'"'+o+">"+i.content+"</"+i.tag+">"},!e.document)return e.addEventListener?(r.disableWorkerMessageHandler||e.addEventListener("message",(function(n){var t=JSON.parse(n.data),a=t.language,i=t.code,l=t.immediateClose;e.postMessage(r.highlight(i,r.languages[a],a)),l&&e.close()}),!1),r):r;var s=r.util.currentScript();function u(){r.manual||r.highlightAll()}if(s&&(r.filename=s.src,s.hasAttribute("data-manual")&&(r.manual=!0)),!r.manual){var c=document.readyState;"loading"===c||"interactive"===c&&s&&s.defer?document.addEventListener("DOMContentLoaded",u):window.requestAnimationFrame?window.requestAnimationFrame(u):window.setTimeout(u,16)}return r}("undefined"!==typeof window?window:"undefined"!==typeof WorkerGlobalScope&&self instanceof WorkerGlobalScope?self:{});e.exports&&(e.exports=t),"undefined"!==typeof n&&(n.Prism=t)}).call(this,t(26))}}]);
//# sourceMappingURL=71.adac62f0.chunk.js.map