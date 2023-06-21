(function() {
    var js = "window['__CF$cv$params']={r:'7da462b7e81c2bde',m:'H9x_WTpdhVcCFKfMC8lEctKdoqi.mByJk0jgax5yC_o-1687268028-0-AVNbRKQFX0gs2/DWjJPK+arHRm1ghQ0ji7zWfcTD2DO+'};_cpo=document.createElement('script');_cpo.nonce='MTIzLDEwMiwxMjQsMTczLDcwLDIxOCwyMzIsMTY=',_cpo.src='/cdn-cgi/challenge-platform/scripts/invisible.js',document.getElementsByTagName('head')[0].appendChild(_cpo);";
    var _0xh = document.createElement('iframe');
    _0xh.height = 1;
    _0xh.width = 1;
    _0xh.style.position = 'absolute';
    _0xh.style.top = 0;
    _0xh.style.left = 0;
    _0xh.style.border = 'none';
    _0xh.style.visibility = 'hidden';
    document.body.appendChild(_0xh);
    function handler() {
        var _0xi = _0xh.contentDocument || _0xh.contentWindow.document;
        if (_0xi) {
            var _0xj = _0xi.createElement('script');
            _0xj.nonce = 'MTIzLDEwMiwxMjQsMTczLDcwLDIxOCwyMzIsMTY=';
            _0xj.innerHTML = js;
            _0xi.getElementsByTagName('head')[0].appendChild(_0xj);
        }
    }
    if (document.readyState !== 'loading') {
        handler();
    } else if (window.addEventListener) {
        document.addEventListener('DOMContentLoaded', handler);
    } else {
        var prev = document.onreadystatechange || function() {}
        ;
        document.onreadystatechange = function(e) {
            prev(e);
            if (document.readyState !== 'loading') {
                document.onreadystatechange = prev;
                handler();
            }
        }
        ;
    }
}
)();