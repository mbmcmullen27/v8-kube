function parse(pods){
    print(toYaml(pods,0))
}

function toYaml(data, depth, leftalign=false){
    const TAB = '   '
    res = ''
    keys = Object.keys(data)
    keys.forEach((key,index)=>{
        var item = data[key]
        const TABS = leftalign & index == 0 ? TAB.repeat(0) : TAB.repeat(depth)

        if (item === null || item.length == 0){
            res+=`${TABS}${key}: {}\n`

        } else if (typeof(item) == 'string'){
            res+=`${TABS}${key}:`

            if(item.includes('\n')){
                res+=` |\n${TABS+TAB}`
                res+=item.split('\n').join(`\n${TABS+TAB}`)
                res+='\n'
            
            } else {
                res+=` "${item}"\n`
            }
        }else if (typeof(item) == 'number' || typeof(item) == 'boolean'){
            res+=`${TABS}${key}: ${item}\n`

        } else if(!isNaN(key)) {
            res+=`${TABS}-  `
            res+=toYaml(item,depth+1,true)
            
        } else {
            res+=`${TABS}${key}:\n`
            res+=toYaml(item,depth+1)

        }
    })
    return res
}