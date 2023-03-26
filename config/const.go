package config

const RBAC_MODEL = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) == true \
    && keyMatch2(r.obj, p.obj) == true \
    && regexMatch(r.act, p.act) == true \
    || r.sub == "9"
`
