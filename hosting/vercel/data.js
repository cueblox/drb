const fs = require('fs');
const path = require('path');

const outputroot = path.join(__dirname, "../.", "output");

const content_types = ["articles", "categories", "pages", "profiles", "sites"]

var data = {}

content_types.forEach((ct) => {

    const files = fs.readdirSync(path.join(outputroot, ct));
    const acc = [];
    files.forEach((curr) => {
        if (curr.includes('.json')) {
            const file = JSON.parse(fs.readFileSync(path.join(outputroot, ct, curr), 'utf8'));
            const extension = path.extname(curr);
            const slug = path.basename(curr, extension);
            // remove wrapping key
            var stripped = file[slug];
            // add id property
            stripped.id = stripped.basename;
            acc.push(stripped);
        }

    }, {});
    data[ct] = acc;


}, {})
module.exports = data;
