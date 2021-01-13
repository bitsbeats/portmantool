CREATE TYPE protocol	AS ENUM ('tcp', 'udp');
CREATE TYPE state	AS ENUM ('open', 'closed', 'filtered', 'unfiltered', 'open|filtered', 'closed|filtered');
