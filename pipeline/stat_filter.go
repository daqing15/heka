/***** BEGIN LICENSE BLOCK *****
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.
#
# The Initial Developer of the Original Code is the Mozilla Foundation.
# Portions created by the Initial Developer are Copyright (C) 2012
# the Initial Developer. All Rights Reserved.
#
# Contributor(s):
#   Ben Bangert (bbangert@mozilla.com)
#   Mike Trinkala (trink@mozilla.com)
#
# ***** END LICENSE BLOCK *****/

package pipeline

type metric struct {
	Type_ string `json:"type"`
	Name  string
	Value string
}

type StatFilterConfig struct {
	Metric []metric
}

type StatFilter struct {
	metrics []metric
}

func (s *StatFilter) ConfigStruct() interface{} {
	return new(StatFilterConfig)
}

func (s *StatFilter) Init(config interface{}) (err error) {
	conf := config.(*StatFilterConfig)
	s.metrics = conf.Metric
	return
}

func (s *StatFilter) Run(fr FilterRunner, h PluginHelper) (err error) {

	inChan := fr.InChan()

	var (
		pack     *PipelinePack
		captures map[string]string
	)
	for plc := range inChan {
		pack = plc.Pack
		captures = plc.Captures

		// Load existing fields into the set for replacement
		captures["Logger"] = pack.Message.GetLogger()
		captures["Hostname"] = pack.Message.GetHostname()
		captures["Type"] = pack.Message.GetType()
		captures["Payload"] = pack.Message.GetPayload()

		// We matched, generate appropriate metrics
		for _, met := range s.metrics {
			m := MessageGenerator.Retrieve()
			m.Message.SetType(met.Type_)
			m.Message.SetLogger(InterpolateString(met.Name, captures))
			m.Message.SetPayload(InterpolateString(met.Value, captures))
			MessageGenerator.Inject(m)
		}
	}

	return
}
