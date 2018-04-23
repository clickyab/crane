/**
 * Clickyab Tracking Script
 *
 * This script attach a function to window object (clickyab_callback)
 * and by calling the this function an impression request will send to clickyab.
 *
 */

;(function() {
  let cyClick: string | null
  let cyImp: string | null

  function getCookie(name: string): string | null {
    const re = new RegExp(name + '=([^;]+)')
    const value = re.exec(document.cookie)
    return value != null ? decodeURI(value[1]) : null
  }

  function decode(s: string) {
    return decodeURIComponent(s.replace(/\+/g, ' '))
  }

  function getQueryStrings() {
    let assoc: {
      [key: string]: string
    } = {};
    const queryString = location.search.substring(1);
    const keyValues = queryString.split('&');
    for (let i in keyValues) {
      let key = keyValues[i].split('=');
      if (key.length > 1) {
        assoc[decode(key[0])] = decode(key[1])
      }
    }
    return assoc
  }

  function appendImgHit(click: string | null, imp: string, actionId: string) {
    const imgHit = document.createElement('img');
    imgHit.setAttribute('src', `{{.URL}}?click_id=${click}&imp_id=${imp}&action_id=${actionId}`);
    document.body.appendChild(imgHit)
  }

  function clickyab_callback(actionId: string) {
    if (actionId === undefined) {
      actionId = ''
    }

    if (cyImp !== null && cyImp !== undefined) {
      appendImgHit(cyClick, cyImp, actionId)
    }
  }

  let getWholeQuery = getQueryStrings();
  cyClick = getWholeQuery['cy_click'];
  cyImp = getWholeQuery['cy_imp'];

  if (cyClick && cyImp) {
    document.cookie = 'cy_click=' + cyClick + '; expires=Fri, 31 Dec 2022 23:59:59 GMT'
    document.cookie = 'cy_imp=' + cyImp + '; expires=Fri, 31 Dec 2022 23:59:59 GMT'
  } else {
    cyClick = getCookie('cy_click');
    cyClick = getCookie('cy_imp');
  }

  ;(window as any).clickyab_callback = clickyab_callback
})()
