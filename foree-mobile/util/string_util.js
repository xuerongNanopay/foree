containSubsequence = (s, q, {caseInsensitive=false}) => {
    if ( typeof s !== 'string' && typeof s !== typeof q) return false
    if ( q.length > s.length ) return false
    if ( caseInsensitive ) {
        s = s.toLowerCase()
        q = q.toLowerCase()
    }

    let p1 = 0;
    let p2 = 0;
    
    while(p1 < q.length) {
      if (q[p1] === s[p2]) {
        p1++;
        p2++;
      } else {
        p2++;
        if (p2 > s.length) {
          return false;
        }
      }
    }
  
    return true;
}

export default {
    containSubsequence
}