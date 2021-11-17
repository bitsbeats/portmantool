CREATE TYPE protocol	AS ENUM ('ip', 'tcp', 'udp', 'sctp');
CREATE TYPE state	AS ENUM ('open', 'closed', 'filtered', 'unfiltered', 'open|filtered', 'closed|filtered');
