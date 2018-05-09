package datastore

import (
	"github.com/spiffe/spire/proto/common"
)

func (d *Dump) AppendBundle(b *Bundle) {
	d.Progress.Dumped.Bundles++
	d.Progress.Dumped.CertBytes += int32(len(b.CaCerts))
	d.Bundles = append(d.Bundles, b)
}

func (d *Dump) AppendJoinToken(j *JoinToken) {
	d.Progress.Dumped.JoinTokens++
	d.JoinTokens = append(d.JoinTokens, j)
}

func (d *Dump) AppendAttestedNodeEntry(a *AttestedNodeEntry) {
	d.Progress.Dumped.AttestedNodeEntries++
	d.AttestedNodeEntries = append(d.AttestedNodeEntries, a)
}

func (d *Dump) AppendNodeResolverMapEntry(n *NodeResolverMapEntry) {
	d.Progress.Dumped.NodeResolverMapEntries++
	d.NodeResolverMapEntries = append(d.NodeResolverMapEntries, n)
}

func (d *Dump) AppendRegistrationEntry(r *common.RegistrationEntry) {
	d.Progress.Dumped.RegistrationEntries++
	d.RegistrationEntries = append(d.RegistrationEntries, r)
}
