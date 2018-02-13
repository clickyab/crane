import CONFIG from './config'
import { IAd } from './interfaces'
import IUrlParams from './interfaces/IUrlParams'
import FP from 'fingerprintjs2'

import Effect from './effect'

declare var clickyabParams: { [key: string]: string }
declare var escape: any

export default class ShowAd {
  private ads: IAd[] = []
  private domain: string
  private publisherId: string

  constructor() {
    if (!clickyabParams.id || !clickyabParams.domain) {
      throw new Error("Clickyab script doesn't config correctly!")
    } else {
      this.domain = clickyabParams.domain
      this.publisherId = clickyabParams.id
    }
  }

  public run() {
    if (window.document.body.getAttribute('clickyab-showAd-ready') === 'true') return
    window.document.body.setAttribute('clickyab-showAd-ready', 'true')
    console.log('start show ad')
    this.ads = this.findAdsInPage()
      .map(elem => this.parseElementProps(elem))
      .map(this.setStyle)
    this.getAdsFromRemote(ads => {
      this.injectSrc(ads)
      this.injectIframes()
      if (ads.m) {
        this.injectMobileAds(ads.m)
      }
    })
  }

  private setStyle(ad: IAd) {
    ad.element.style.height = ad.height + 'px'
    ad.element.style.textAlign = 'center'
    return ad
  }

  private injectIframes() {
    this.ads.forEach(ad => {
      let ignoreAdBecauseCookie = false
      if (ad.valid) {
        if (ad.effect === 'interstitial' && this.getCookie('cy_interstitial')) {
          ignoreAdBecauseCookie = true
        } else if (ad.effect === 'interstitial' && !this.getCookie('cy_interstitial')) {
          this.setCookie('cy_interstitial', 'true', 0.5)
        }

        if (!ignoreAdBecauseCookie && ad.iframe) {
          ad.element.innerHTML = decodeURIComponent(ad.iframe || '')

          if (ad.effect) {
            const effectAct = new Effect(ad)
          }
        }
      }
    })
  }

  private injectSrc(ads: { [key: string]: string }) {
    this.ads = this.ads.map(ad => {
      ad.iframe = ads[`${ad.slot}`]
      return ad
    })
  }

  private getAdsFromRemote(onload: (ads: { [key: string]: string }) => void) {
    this.generateUrl(url => {
      let request = new XMLHttpRequest()
      request.addEventListener('load', function() {
        try {
          onload(JSON.parse(this.responseText))
        } catch (err) {
          console.log('Error in get ads list.')
        }
      })

      request.open('GET', url)
      request.send()
    })
  }

  private generateUrl(onLoad: (url: string) => void) {
    let url: IUrlParams = {
      tracking: 'true',
      mobile: this.isMobile(),
      id: this.publisherId,
      domain: this.domain,
      loc: this.encodeuri(document.location.href),
      ref: this.encodeuri(document.referrer),
      count: this.ads.length.toString(),
      slots: this.ads
        .filter(a => a.valid)
        .map(a => {
          return `${a.slot}:${a.size}`
        })
        .join(',')
    }

    let urlParamsString: string[] = Object.keys(url).map(s => {
      return `${s[0]}=${url[s].toString()}`
    })

    new FP().get((r: string) => {
      urlParamsString.push(`tid=${r}`)

      let UrlString = `${CONFIG.REMOTE_TARGET}?${urlParamsString.join('&')}`
      onLoad(UrlString)
    })
  }

  private isMobile(): number {
    if (
      navigator.userAgent.match(/Android/i) ||
      navigator.userAgent.match(/webOS/i) ||
      navigator.userAgent.match(/iPhone/i) ||
      navigator.userAgent.match(/iPad/i) ||
      navigator.userAgent.match(/iPod/i) ||
      navigator.userAgent.match(/BlackBerry/i) ||
      navigator.userAgent.match(/Windows Phone/i)
    ) {
      return 1
    } else {
      return 0
    }
  }

  private findAdsInPage(): HTMLElement[] {
    let elements: HTMLElement[] = []
    const elementsCollection = document.getElementsByClassName(CONFIG.SELECTOR_CLASS)
    for (let i = 0; i < elementsCollection.length; i++) {
      elements.push(elementsCollection.item(i) as HTMLElement)
    }
    return elements
  }

  private parseElementProps(element: HTMLElement): IAd {
    let ad: IAd = {
      element: element,
      slot: element.getAttribute(CONFIG.ELEMENT_PROPERTY_PREFIX + 'slot'),
      height: element.getAttribute(CONFIG.ELEMENT_PROPERTY_PREFIX + 'height'),
      width: element.getAttribute(CONFIG.ELEMENT_PROPERTY_PREFIX + 'width'),
      minFlex: element.getAttribute(CONFIG.ELEMENT_PROPERTY_PREFIX + 'minFlex'),
      effect: element.getAttribute(CONFIG.ELEMENT_PROPERTY_PREFIX + 'effect')
    }
    ad.valid = this.validateAdElement(ad)
    ad.size = this.getAdSize(ad)
    return ad
  }

  private validateAdElement(ad: IAd): boolean {
    return this.getAdSize(ad) !== -1
  }

  private injectMobileAds(src: string) {
    const div = document.createElement('div');
    div.setAttribute(
      'style',
      `position: fixed; width: 100%; z-index:99999999; left: 0; bottom: 0px; margin: 0; padding: 0; text-align: center;`
    );
    div.style.height = "50px";
    div.innerHTML = src;

    document.getElementsByTagName('body')[0].appendChild(div)
  }

  private getAdSize(ad: IAd): number {
    const size: string = `${ad.width}_${ad.height}`
    const sizes: { [index: string]: number } = CONFIG.BANNER_SIZES as { [index: string]: number }
    return sizes[size] || -1
  }

  private encodeuri(b: string): string {
    if (typeof encodeURIComponent === 'function') {
      return encodeURIComponent(b)
    } else {
      return escape(b)
    }
  }

  private setCookie(cname: string, cvalue: string, exdays: number) {
    const d = new Date()
    d.setTime(d.getTime() + exdays * 24 * 60 * 60 * 1000)
    const expires = 'expires=' + d.toUTCString()
    document.cookie = cname + '=' + cvalue + '; ' + expires
  }

  private getCookie(cname: string): string {
    const name = cname + '='
    const ca = document.cookie.split(';')
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i]
      while (c.charAt(0) === ' ') {
        c = c.substring(1)
      }
      if (c.indexOf(name) === 0) {
        return c.substring(name.length, c.length)
      }
    }
    return ''
  }
}
